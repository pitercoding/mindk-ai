package routes

import (
	"net/http"
	"time"

	"github.com/pitercoding/mindk-ai/backend/internal/app"
	"github.com/pitercoding/mindk-ai/backend/internal/handlers"
	"github.com/pitercoding/mindk-ai/backend/internal/middleware"
)

// Rate limits are per authenticated user. Chat is limited more strictly
// than the rest of the API because it drives RAG lookups and LLM calls,
// which are far more expensive than a CRUD request.
const (
	generalRateLimit = 60
	chatRateLimit    = 10
	rateLimitWindow  = time.Minute
)

func RegisterRoutes(app *app.App) {

	protected := app.AuthMiddleware
	general := middleware.NewRateLimiter(generalRateLimit, rateLimitWindow).Middleware
	chatLimit := middleware.NewRateLimiter(chatRateLimit, rateLimitWindow).Middleware

	// Public routes
	http.HandleFunc("/health", handlers.HealthHandler)

	// Protected routes
	http.Handle("/notes", protected(general(http.HandlerFunc(app.NoteHandler.HandleNotes))))
	http.Handle("/notes/", protected(general(http.HandlerFunc(app.NoteHandler.HandleNote))))

	http.Handle("/chat", protected(chatLimit(http.HandlerFunc(app.ChatHandler.Ask))))

	http.Handle("/chat/messages/", protected(general(http.HandlerFunc(app.ChatMessageHandler.HandleMessages))))

	http.Handle("/documents", protected(general(http.HandlerFunc(app.DocumentHandler.HandleDocuments))))
	http.Handle("/documents/upload", protected(general(http.HandlerFunc(app.DocumentHandler.UploadDocument))))
	http.Handle("/documents/search", protected(general(http.HandlerFunc(app.DocumentHandler.SearchDocuments))))
	http.Handle("/documents/", protected(general(http.HandlerFunc(app.DocumentHandler.HandleDocument))))

	http.Handle("/chat/sessions", protected(general(http.HandlerFunc(app.ChatSessionHandler.HandleSessions))))
	http.Handle("/chat/sessions/", protected(general(http.HandlerFunc(app.ChatSessionHandler.HandleSession))))
}
