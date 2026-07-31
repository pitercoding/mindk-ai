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
	ChatMessageHandler *handlers.ChatMessageHandler
	DocumentHandler    *handlers.DocumentHandler
}

func New(
	db *sql.DB,
	cfg *config.Config,
) *App {

	// Repository
	noteRepo := repository.NewNoteRepository(db)
	chatMessageRepo := repository.NewChatMessageRepository(db)
	documentRepo := repository.NewDocumentRepository(db)
	documentChunkRepo := repository.NewDocumentChunkRepository(db)

	// Services
	noteService := services.NewNoteService(noteRepo)
	chatMessageService := services.NewChatMessageService(chatMessageRepo)
	documentChunkService := services.NewDocumentChunkService(documentChunkRepo)

	documentService := services.NewDocumentService(
		documentRepo, 
		documentChunkService,
	)

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
	documentHandler := handlers.NewDocumentHandler(documentService)

	return &App{
		NoteHandler:        noteHandler,
		ChatHandler:        chatHandler,
		ChatMessageHandler: chatMessageHandler,
		DocumentHandler:    documentHandler,
	}
}
