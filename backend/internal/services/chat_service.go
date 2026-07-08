package services

import (
	"github.com/pitercoding/mindk-ai/backend/internal/llm"
)

type ChatService struct {
	noteService        *NoteService
	chatHistoryService *ChatHistoryService
	llmClient          llm.Client
}

func NewChatService(
	noteService *NoteService,
	chatHistoryService *ChatHistoryService,
	llmClient llm.Client,
) *ChatService {

	return &ChatService{
		noteService:        noteService,
		chatHistoryService: chatHistoryService,
		llmClient:          llmClient,
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

	err = s.chatHistoryService.Save(message, answer)
	if err != nil {
		return "", err
	}

	return answer, nil
}
