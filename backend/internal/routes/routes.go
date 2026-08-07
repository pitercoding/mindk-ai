package routes

import (
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/app"
	"github.com/pitercoding/mindk-ai/backend/internal/handlers"
)

func RegisterRoutes(app *app.App) {

	// Routes
	http.HandleFunc("/health", handlers.HealthHandler)

	http.HandleFunc("/notes", app.NoteHandler.HandleNotes)
	http.HandleFunc("/notes/", app.NoteHandler.HandleNote)

	http.HandleFunc("/chat", app.ChatHandler.Ask)

	http.HandleFunc("/chat/messages/", app.ChatMessageHandler.HandleMessages)

	http.HandleFunc("/documents", app.DocumentHandler.HandleDocuments)
	http.HandleFunc("/documents/upload", app.DocumentHandler.UploadDocument)
	http.HandleFunc("/documents/search", app.DocumentHandler.SearchDocuments)
	http.HandleFunc("/documents/", app.DocumentHandler.HandleDocument)

	http.HandleFunc("/chat/sessions", app.ChatSessionHandler.HandleSessions)
	http.HandleFunc("/chat/sessions/", app.ChatSessionHandler.HandleSession)
}
