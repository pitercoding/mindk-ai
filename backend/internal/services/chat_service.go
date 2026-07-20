package services

import (
	"github.com/pitercoding/mindk-ai/backend/internal/llm"
	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type NoteProvider interface {
	GetAll() ([]models.Note, error)
}

type ChatHistoryProvider interface {
	Save(question, answer string) error
	GetRecent(limit int) ([]models.ChatHistory, error)
}

type ChatService struct {
	noteService        NoteProvider
	chatHistoryService ChatHistoryProvider
	llmClient          llm.Client
}

func NewChatService(
	noteService NoteProvider,
	chatHistoryService ChatHistoryProvider,
	llmClient llm.Client,
) *ChatService {

	return &ChatService{
		noteService:        noteService,
		chatHistoryService: chatHistoryService,
		llmClient:          llmClient,
	}
}

func (s *ChatService) Ask(message string, context *models.ChatContext) (string, error) {

	history, err := s.chatHistoryService.GetRecent(5)
	if err != nil {
		return "", err
	}

	var notes []models.Note

	if context != nil {

		notes = []models.Note{
			{
				Title:   context.Title,
				Content: context.Content,
			},
		}

	} else {

		notes, err = s.noteService.GetAll()

		if err != nil {
			return "", err
		}
	}

	prompt := llm.BuildPrompt(
		message,
		notes,
		history,
	)

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
