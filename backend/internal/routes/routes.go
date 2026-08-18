package routes

import (
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/app"
	"github.com/pitercoding/mindk-ai/backend/internal/handlers"
)

func RegisterRoutes(app *app.App) {

	protected := app.AuthMiddleware

	// Public routes
	http.HandleFunc("/health", handlers.HealthHandler)

	// Protected routes
	http.Handle("/notes", protected(http.HandlerFunc(app.NoteHandler.HandleNotes)))
	http.Handle("/notes/", protected(http.HandlerFunc(app.NoteHandler.HandleNote)))

	http.Handle("/chat", protected(http.HandlerFunc(app.ChatHandler.Ask)))

	http.Handle("/chat/messages/", protected(http.HandlerFunc(app.ChatMessageHandler.HandleMessages)))

	http.Handle("/documents", protected(http.HandlerFunc(app.DocumentHandler.HandleDocuments)))
	http.Handle("/documents/upload", protected(http.HandlerFunc(app.DocumentHandler.UploadDocument)))
	http.Handle("/documents/search", protected(http.HandlerFunc(app.DocumentHandler.SearchDocuments)))
	http.Handle("/documents/", protected(http.HandlerFunc(app.DocumentHandler.HandleDocument)))

	http.Handle("/chat/sessions", protected(http.HandlerFunc(app.ChatSessionHandler.HandleSessions)))
	http.Handle("/chat/sessions/", protected(http.HandlerFunc(app.ChatSessionHandler.HandleSession)))
}
