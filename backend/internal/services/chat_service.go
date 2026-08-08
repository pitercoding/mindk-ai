package services

import (
	"errors"

	"github.com/pitercoding/mindk-ai/backend/internal/llm"
	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type NoteProvider interface {
	GetAll() ([]models.Note, error)
	GetByID(id int) (*models.Note, error)
}

type ChatSessionProvider interface {
	GetByID(id int) (*models.ChatSession, error)
}

type ChatMessageProvider interface {
	Save(message *models.ChatMessage) error
	GetBySessionID(sessionID int) ([]models.ChatMessage, error)
}

type DocumentContextProvider interface {
	BuildContext(
		query string,
		limit int,
	) (string, error)
}

type PromptBuilder interface {
	Build(
		question string,
		notes []models.Note,
		messages []models.ChatMessage,
		documentContext string,
	) string
}

type ChatService struct {
	noteService        NoteProvider
	chatSessionService ChatSessionProvider
	documentContext    DocumentContextProvider
	chatMessageService ChatMessageProvider
	llmClient          llm.Client
	promptBuilder      PromptBuilder
}

func NewChatService(
	noteService NoteProvider,
	chatSessionService ChatSessionProvider,
	documentContext DocumentContextProvider,
	chatMessageService ChatMessageProvider,
	llmClient llm.Client,
) *ChatService {

	return &ChatService{
		noteService:        noteService,
		chatSessionService: chatSessionService,
		documentContext:    documentContext,
		chatMessageService: chatMessageService,
		llmClient:          llmClient,
		promptBuilder:      llm.NewContextBuilder(),
	}
}

func (s *ChatService) Ask(
	sessionID int,
	message string,
) (string, error) {

	session, err := s.chatSessionService.GetByID(sessionID)

	if err != nil {
		return "", err
	}

	if session == nil {
		return "", errors.New("chat session not found")
	}

	messages, err := s.chatMessageService.GetBySessionID(
		sessionID,
	)

	if err != nil {
		return "", err
	}

	var notes []models.Note

	switch session.Mode {

	case "knowledge":

		notes, err = s.noteService.GetAll()

		if err != nil {
			return "", err
		}

	case "note":

		if session.NoteID == nil {
			return "", errors.New(
				"note session has no note",
			)
		}

		note, err := s.noteService.GetByID(
			*session.NoteID,
		)

		if err != nil {
			return "", err
		}

		if note == nil {
			return "", errors.New("note not found")
		}

		notes = []models.Note{
			*note,
		}

	default:

		return "", errors.New(
			"invalid chat session mode",
		)
	}

	documentContext := ""

	if session.Mode == "knowledge" &&
		s.documentContext != nil {

		documentContext, err =
			s.documentContext.BuildContext(
				message,
				5,
			)

		if err != nil {
			return "", err
		}
	}

	prompt := s.promptBuilder.Build(
		message,
		notes,
		messages,
		documentContext,
	)

	answer, err := s.llmClient.Chat(prompt)

	if err != nil {
		return "", err
	}

	err = s.chatMessageService.Save(
		&models.ChatMessage{
			SessionID: sessionID,
			Role:      "user",
			Content:   message,
		},
	)

	if err != nil {
		return "", err
	}

	err = s.chatMessageService.Save(
		&models.ChatMessage{
			SessionID: sessionID,
			Role:      "assistant",
			Content:   answer,
		},
	)

	if err != nil {
		return "", err
	}

	return answer, nil
}
