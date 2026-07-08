package app

import (
	"database/sql"

	"github.com/pitercoding/mindk-ai/backend/internal/config"
	"github.com/pitercoding/mindk-ai/backend/internal/handlers"
	"github.com/pitercoding/mindk-ai/backend/internal/llm"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
	"github.com/pitercoding/mindk-ai/backend/internal/services"
)

type App struct {
	NoteHandler *handlers.NoteHandler
	ChatHandler *handlers.ChatHandler
	ChatHistoryHandler *handlers.ChatHistoryHandler
}

func New(
	db *sql.DB,
	cfg *config.Config,
) *App {

	// Repository
	noteRepo := repository.NewNoteRepository(db)
	chatRepo := repository.NewChatRepository(db)

	// Services
	noteService := services.NewNoteService(noteRepo)
	chatHistoryService := services.NewChatHistoryService(chatRepo)

	// LLM Client
	openAIClient := llm.NewOpenAIClient(
		cfg.OpenAIAPIKey,
	)

	// Chat Service
	chatService := services.NewChatService(
		noteService,
		chatHistoryService,
		openAIClient,
	)

	// Handlers
	noteHandler := handlers.NewNoteHandler(noteService)
	chatHandler := handlers.NewChatHandler(chatService)
	chatHistoryHandler := handlers.NewChatHistoryHandler(chatHistoryService,)

	return &App{
		NoteHandler: noteHandler,
		ChatHandler: chatHandler,
		ChatHistoryHandler: chatHistoryHandler,
	}
}
