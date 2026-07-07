package services

import (
	"github.com/pitercoding/mindk-ai/backend/internal/llm"
)

type ChatService struct {
	noteService *NoteService
	llmClient   llm.Client
}

func NewChatService(
	noteService *NoteService,
	llmClient llm.Client,
) *ChatService {

	return &ChatService{
		noteService: noteService,
		llmClient:   llmClient,
	}
}

func (s *ChatService) Ask(message string) (string, error) {
	notes, err := s.noteService.GetAll()

	if err != nil {
		return "", err
	}

	prompt := llm.BuildPrompt(message, notes)

	answer, err := s.llmClient.Chat(prompt)

	if err != nil {
		return "", err
	}

	return answer, nil
}
