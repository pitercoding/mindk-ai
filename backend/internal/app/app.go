package app

import (
	"database/sql"

	"github.com/pitercoding/mindk-ai/backend/internal/handlers"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
	"github.com/pitercoding/mindk-ai/backend/internal/services"
)

type App struct {
	NoteHandler *handlers.NoteHandler
	ChatHandler *handlers.ChatHandler
}

func New(db *sql.DB) *App {
	noteRepo := repository.NewNoteRepository(db)

	noteService := services.NewNoteService(noteRepo)
	chatService := services.NewChatService(noteService)

	noteHandler := handlers.NewNoteHandler(noteService)
	chatHandler := handlers.NewChatHandler(chatService)

	return &App{
		NoteHandler: noteHandler,
		ChatHandler: chatHandler,
	}
}
