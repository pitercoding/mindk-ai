package routes

import (
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/pitercoding/mindk-ai/backend/docs"
	"github.com/pitercoding/mindk-ai/backend/internal/app"
	"github.com/pitercoding/mindk-ai/backend/internal/config"
	"github.com/pitercoding/mindk-ai/backend/internal/handlers"
	"github.com/pitercoding/mindk-ai/backend/internal/middleware"
)

// Rate limits are per authenticated user. Chat is limited more strictly
// than the rest of the API because it drives RAG lookups and LLM calls,
// which are far more expensive than a CRUD request. Exported so tests can
// assert against the real production limits instead of duplicating them.
const (
	GeneralRateLimit = 60
	ChatRateLimit    = 10
	RateLimitWindow  = time.Minute
)

// RegisterRoutes wires every handler in app onto mux, exactly as production
// serves them. Taking mux explicitly (rather than registering on
// http.DefaultServeMux) lets each test build its own isolated server without
// route registrations leaking between tests.
func RegisterRoutes(mux *http.ServeMux, app *app.App) {

	protected := app.AuthMiddleware
	general := middleware.NewRateLimiter(GeneralRateLimit, RateLimitWindow).Middleware
	chatLimit := middleware.NewRateLimiter(ChatRateLimit, RateLimitWindow).Middleware

	// Public routes
	mux.HandleFunc("/health", handlers.HealthHandler)

	// Protected routes
	mux.Handle("/notes", protected(general(http.HandlerFunc(app.NoteHandler.HandleNotes))))
	mux.Handle("/notes/", protected(general(http.HandlerFunc(app.NoteHandler.HandleNote))))

	mux.Handle("/chat", protected(chatLimit(http.HandlerFunc(app.ChatHandler.Ask))))

	mux.Handle("/chat/messages/", protected(general(http.HandlerFunc(app.ChatMessageHandler.HandleMessages))))

	mux.Handle("/documents", protected(general(http.HandlerFunc(app.DocumentHandler.HandleDocuments))))
	mux.Handle("/documents/upload", protected(general(http.HandlerFunc(app.DocumentHandler.UploadDocument))))
	mux.Handle("/documents/search", protected(general(http.HandlerFunc(app.DocumentHandler.SearchDocuments))))
	mux.Handle("/documents/", protected(general(http.HandlerFunc(app.DocumentHandler.HandleDocument))))

	mux.Handle("/chat/sessions", protected(general(http.HandlerFunc(app.ChatSessionHandler.HandleSessions))))
	mux.Handle("/chat/sessions/", protected(general(http.HandlerFunc(app.ChatSessionHandler.HandleSession))))

	// Swagger UI is a development aid only. It is never registered in
	// production, so /swagger/ simply 404s there instead of needing its own
	// auth story.
	if app.Environment == config.EnvDevelopment {
		mux.Handle("/swagger/", httpSwagger.WrapHandler)
	}
}
