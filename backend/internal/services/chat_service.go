package services

import (
	"github.com/pitercoding/mindk-ai/backend/internal/llm"
	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type NoteProvider interface {
	GetAll() ([]models.Note, error)
}

type ChatMessageProvider interface {
	Save(message *models.ChatMessage) error
	GetByNoteID(noteID int) ([]models.ChatMessage, error)
}

type ChatService struct {
	noteService        NoteProvider
	chatMessageService ChatMessageProvider
	llmClient          llm.Client
}

func NewChatService(
	noteService NoteProvider,
	chatMessageService ChatMessageProvider,
	llmClient llm.Client,
) *ChatService {

	return &ChatService{
		noteService:        noteService,
		chatMessageService: chatMessageService,
		llmClient:          llmClient,
	}
}

func (s *ChatService) Ask(
	message string,
	context *models.ChatContext,
) (string, error) {

	var (
		err      error
		notes    []models.Note
		messages []models.ChatMessage
	)

	if context != nil {

		notes = []models.Note{
			{
				Title:   context.Title,
				Content: context.Content,
			},
		}

		messages, err = s.chatMessageService.GetByNoteID(
			context.NoteID,
		)

		if err != nil {
			return "", err
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
		messages,
	)

	answer, err := s.llmClient.Chat(prompt)

	if err != nil {
		return "", err
	}

	if context != nil {

		err = s.chatMessageService.Save(
			&models.ChatMessage{
				NoteID:  context.NoteID,
				Role:    "user",
				Content: message,
			},
		)

		if err != nil {
			return "", err
		}

		err = s.chatMessageService.Save(
			&models.ChatMessage{
				NoteID:  context.NoteID,
				Role:    "assistant",
				Content: answer,
			},
		)

		if err != nil {
			return "", err
		}
	}

	return answer, nil
}
