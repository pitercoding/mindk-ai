package services

import (
	"errors"

	"github.com/pitercoding/mindk-ai/backend/internal/llm"
	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
)

type NoteProvider interface {
	GetAll() ([]models.Note, error)
	GetByID(id int) (*models.Note, error)
}

type ChatSessionProvider interface {
	GetByID(id int) (*models.ChatSession, error)
	Create(session *models.ChatSession) error
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

// Ask resolves (or creates) the chat session, persists the user's message,
// builds the prompt from notes/documents/history and returns the
// assistant's answer together with the session it belongs to.
func (s *ChatService) Ask(
	sessionID int,
	message string,
	mode string,
	noteID *int,
	title string,
) (string, int, error) {

	session, err := s.resolveSession(
		sessionID,
		mode,
		noteID,
		title,
	)

	if err != nil {
		return "", 0, err
	}

	// History must be fetched before saving the current message,
	// otherwise the question would be duplicated in the prompt
	// (once as history, once as the question itself).
	history, err := s.chatMessageService.GetBySessionID(
		session.ID,
	)

	if err != nil {
		return "", 0, err
	}

	// Saved before calling the LLM so the user's input is never lost
	// if the request to OpenAI fails afterwards.
	err = s.chatMessageService.Save(
		&models.ChatMessage{
			SessionID: session.ID,
			Role:      "user",
			Content:   message,
		},
	)

	if err != nil {
		return "", 0, err
	}

	var notes []models.Note

	switch session.Mode {

	case "knowledge":

		notes, err = s.noteService.GetAll()

		if err != nil {
			return "", 0, err
		}

	case "note":

		if session.NoteID == nil {
			return "", 0, errors.New(
				"note session has no note",
			)
		}

		note, err := s.noteService.GetByID(
			*session.NoteID,
		)

		if err != nil {
			return "", 0, err
		}

		if note == nil {
			return "", 0, errors.New("note not found")
		}

		notes = []models.Note{
			*note,
		}

	default:

		return "", 0, errors.New(
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
			return "", 0, err
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
		return "", 0, err
	}

	err = s.chatMessageService.Save(
		&models.ChatMessage{
			SessionID: session.ID,
			Role:      "assistant",
			Content:   answer,
		},
	)

	if err != nil {
		return "", 0, err
	}

	return answer, session.ID, nil
}

// resolveSession looks up an existing session, or creates a new one when
// sessionID is not provided (0). Auto-creation requires a mode, since a
// chat session cannot be built without knowing how to answer questions.
func (s *ChatService) resolveSession(
	sessionID int,
	mode string,
	noteID *int,
	title string,
) (*models.ChatSession, error) {

	if sessionID != 0 {

		session, err := s.chatSessionService.GetByID(
			sessionID,
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
