package app

import (
	"database/sql"
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/config"
	"github.com/pitercoding/mindk-ai/backend/internal/handlers"
	"github.com/pitercoding/mindk-ai/backend/internal/llm"
	"github.com/pitercoding/mindk-ai/backend/internal/middleware"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
	"github.com/pitercoding/mindk-ai/backend/internal/services"
)

type App struct {
	NoteHandler        *handlers.NoteHandler
	ChatHandler        *handlers.ChatHandler
	ChatMessageHandler *handlers.ChatMessageHandler
	ChatSessionHandler *handlers.ChatSessionHandler
	DocumentHandler    *handlers.DocumentHandler
	AuthMiddleware     func(http.Handler) http.Handler
}

func New(
	db *sql.DB,
	cfg *config.Config,
) *App {

	// Repository
	noteRepo := repository.NewNoteRepository(db)
	chatMessageRepo := repository.NewChatMessageRepository(db)
	chatSessionRepo := repository.NewChatSessionRepository(db)
	documentRepo := repository.NewDocumentRepository(db)
	documentChunkRepo := repository.NewDocumentChunkRepository(db)
	documentEmbeddingRepo := repository.NewDocumentEmbeddingRepository(db)

	// LLM Client
	openAIClient := llm.NewOpenAIClient(cfg.OpenAIAPIKey)
	embeddingGenerator := llm.NewOpenAIEmbeddingGenerator(openAIClient)

	// Services
	noteService := services.NewNoteService(noteRepo)

	chatMessageService := services.NewChatMessageService(
		chatMessageRepo,
		chatSessionRepo,
	)

	chatSessionService := services.NewChatSessionService(
		chatSessionRepo,
	)

	documentChunkService := services.NewDocumentChunkService(
		documentChunkRepo,
	)

	documentEmbeddingService := services.NewDocumentEmbeddingService(
		documentEmbeddingRepo,
		embeddingGenerator,
	)

	documentSearchService := services.NewDocumentSearchService(
		documentEmbeddingRepo,
		embeddingGenerator,
		services.DefaultMinRelevanceScore,
	)

	documentContextService := services.NewDocumentContextService(documentSearchService)

	documentService := services.NewDocumentService(
		documentRepo,
		documentChunkService,
		documentEmbeddingService,
	)

	// Chat Service
	chatService := services.NewChatService(
		noteService,
		chatSessionService,
		documentContextService,
		chatMessageService,
		openAIClient,
	)

	// Handlers
	noteHandler := handlers.NewNoteHandler(noteService)
	chatHandler := handlers.NewChatHandler(chatService)
	chatMessageHandler := handlers.NewChatMessageHandler(chatMessageService)
	chatSessionHandler := handlers.NewChatSessionHandler(chatSessionService)
	documentHandler := handlers.NewDocumentHandler(documentService)

	// Middleware
	authMiddleware := middleware.NewClerkAuth(cfg.ClerkSecretKey)

	return &App{
		NoteHandler:        noteHandler,
		ChatHandler:        chatHandler,
		ChatMessageHandler: chatMessageHandler,
		ChatSessionHandler: chatSessionHandler,
		DocumentHandler:    documentHandler,
		AuthMiddleware:     authMiddleware,
	}
}
