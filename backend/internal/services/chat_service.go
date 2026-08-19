package services

import (
	"errors"

	"github.com/pitercoding/mindk-ai/backend/internal/llm"
	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
)

type NoteProvider interface {
	GetAll(userID string) ([]models.Note, error)
	GetByID(id int, userID string) (*models.Note, error)
}

type ChatSessionProvider interface {
	GetByID(id int, userID string) (*models.ChatSession, error)
	Create(session *models.ChatSession) error
}

type ChatMessageProvider interface {
	Save(message *models.ChatMessage, userID string) error
	GetBySessionID(sessionID int, userID string) ([]models.ChatMessage, error)
}

type DocumentContextProvider interface {
	BuildContext(
		query string,
		limit int,
		userID string,
	) (string, []models.ChatSource, error)
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

// Ask resolves (or creates) the chat session, persists the user's message, builds the prompt from notes/documents/history and returns the assistant's answer together with the session it belongs to.
func (s *ChatService) Ask(
	userID string,
	sessionID int,
	message string,
	mode string,
	noteID *int,
	title string,
) (string, int, []models.ChatSource, error) {

	session, err := s.resolveSession(
		userID,
		sessionID,
		mode,
		noteID,
		title,
	)

	if err != nil {
		return "", 0, nil, err
	}

	// History must be fetched before saving the current message, otherwise the question would be duplicated in the prompt (once as history, once as the question itself).
	history, err := s.chatMessageService.GetBySessionID(
		session.ID,
		userID,
	)

	if err != nil {
		return "", 0, nil, err
	}

	// Saved before calling the LLM so the user's input is never lost if the request to OpenAI fails afterwards.
	err = s.chatMessageService.Save(
		&models.ChatMessage{
			SessionID: session.ID,
			Role:      "user",
			Content:   message,
		},
		userID,
	)

	if err != nil {
		return "", 0, nil, err
	}

	var notes []models.Note

	switch session.Mode {

	case "knowledge":

		notes, err = s.noteService.GetAll(userID)

		if err != nil {
			return "", 0, nil, err
		}

	case "note":

		if session.NoteID == nil {
			return "", 0, nil, errors.New(
				"note session has no note",
			)
		}

		note, err := s.noteService.GetByID(
			*session.NoteID,
			userID,
		)

		if err != nil {
			return "", 0, nil, err
		}

		if note == nil {
			return "", 0, nil, errors.New("note not found")
		}

		notes = []models.Note{
			*note,
		}

	default:

		return "", 0, nil, errors.New(
			"invalid chat session mode",
		)
	}

	documentContext := ""
	var documentSources []models.ChatSource

	if session.Mode == "knowledge" &&
		s.documentContext != nil {

		documentContext, documentSources, err =
			s.documentContext.BuildContext(
				message,
				5,
				userID,
			)

		if err != nil {
			return "", 0, nil, err
		}
	}

	prompt := s.promptBuilder.Build(
		message,
		notes,
		history,
		documentContext,
	)

	answer, err := s.llmClient.Chat(prompt)

	if err != nil {
		return "", 0, nil, err
	}

	err = s.chatMessageService.Save(
		&models.ChatMessage{
			SessionID: session.ID,
			Role:      "assistant",
			Content:   answer,
		},
		userID,
	)

	if err != nil {
		return "", 0, nil, err
	}

	return answer, session.ID, documentSources, nil
}

// resolveSession looks up an existing session, or creates a new one when sessionID is not provided (0). Auto-creation requires a mode, since a chat session cannot be built without knowing how to answer questions.
func (s *ChatService) resolveSession(
	userID string,
	sessionID int,
	mode string,
	noteID *int,
	title string,
) (*models.ChatSession, error) {

	if sessionID != 0 {

		session, err := s.chatSessionService.GetByID(
			sessionID,
			userID,
		)

		if err != nil {
			return nil, err
		}

		if session == nil {
			return nil, repository.ErrChatSessionNotFound
		}

		return session, nil
	}

	if mode == "" {
		return nil, errors.New(
			"mode is required to start a new chat session",
		)
	}

	session := &models.ChatSession{
		UserID: userID,
		Title:  sessionTitle(title),
		Mode:   mode,
		NoteID: noteID,
	}

	err := s.chatSessionService.Create(session)

	if err != nil {
		return nil, err
	}

	return session, nil
}

func sessionTitle(title string) string {

	if title == "" {
		return "New Chat"
	}

	return title
}
