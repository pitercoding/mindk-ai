package app

import (
	"database/sql"

	"github.com/pitercoding/mindk-ai/backend/internal/handlers"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
	"github.com/pitercoding/mindk-ai/backend/internal/services"
)

type App struct {
	NoteHandler *handlers.NoteHandler
}

func New(db *sql.DB) *App {
	noteRepo := repository.NewNoteRepository(db)

	noteService := services.NewNoteService(noteRepo)

	noteHandler := handlers.NewNoteHandler(noteService)

	return &App{
		NoteHandler: noteHandler,
	}
}
