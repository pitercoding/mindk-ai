package app

import (
	"database/sql"

	"github.com/pitercoding/mindk-ai/backend/internal/handlers"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
)

type App struct {
	NoteHandler *handlers.NoteHandler
}

func New(db *sql.DB) *App {
	noteRepo := repository.NewNoteRepository(db)

	noteHandler := handlers.NewNoteHandler(noteRepo)

	return &App{
		NoteHandler: noteHandler,
	}
}
