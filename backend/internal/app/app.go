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
	NoteHandler        *handlers.NoteHandler
	ChatHandler        *handlers.ChatHandler
	ChatHistoryHandler *handlers.ChatHistoryHandler
	ChatMessageHandler *handlers.ChatMessageHandler
}

func New(
	db *sql.DB,
	cfg *config.Config,
) *App {

	// Repository
	noteRepo := repository.NewNoteRepository(db)
	chatMessageRepo := repository.NewChatMessageRepository(db)

	// Services
	noteService := services.NewNoteService(noteRepo)
	chatMessageService := services.NewChatMessageService(chatMessageRepo)

	// LLM Client
	openAIClient := llm.NewOpenAIClient(
		cfg.OpenAIAPIKey,
	)

	// Chat Service
	chatService := services.NewChatService(
		noteService,
		chatMessageService,
		openAIClient,
	)

	// Handlers
	noteHandler := handlers.NewNoteHandler(noteService)
	chatHandler := handlers.NewChatHandler(chatService)
	chatMessageHandler := handlers.NewChatMessageHandler(chatMessageService)

	return &App{
		NoteHandler:        noteHandler,
		ChatHandler:        chatHandler,
		ChatMessageHandler: chatMessageHandler,
	}
}
