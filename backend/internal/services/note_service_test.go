package services

import (
	"errors"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoteServiceGetAll(t *testing.T) {

	tests := []struct {
		name          string
		repoNotes     []models.Note
		repoErr       error
		expectedNotes []models.Note
		expectedCount int
		expectError   bool
	}{
		{
			name: "returns all notes",
			repoNotes: []models.Note{
				{
					ID:      1,
					Title:   "Go",
					Content: "Go is awesome",
				},
				{
					ID:      2,
					Title:   "Docker",
					Content: "Docker containers",
				},
			},
			expectedNotes: []models.Note{
				{
					ID:      1,
					Title:   "Go",
					Content: "Go is awesome",
				},
				{
					ID:      2,
					Title:   "Docker",
					Content: "Docker containers",
				},
			},
			expectedCount: 2,
		},
		{
			name:          "returns empty list",
			repoNotes:     []models.Note{},
			expectedNotes: []models.Note{},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:          "repository returns error",
			repoErr:       errors.New("database error"),
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo := &mocks.FakeNoteRepository{
				Notes: tt.repoNotes,
				Err:   tt.repoErr,
			}

			service := NewNoteService(repo)

			notes, err := service.GetAll()

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, notes, tt.expectedCount)
			assert.Equal(t, tt.expectedNotes, notes)
		})
	}
}

func TestNoteServiceGetByID(t *testing.T) {

	tests := []struct {
		name        string
		repoNote    *models.Note
		repoErr     error
		expectedID  int
		expectedErr bool
	}{
		{
			name: "returns note by id",
			repoNote: &models.Note{
				ID:      1,
				Title:   "Go",
				Content: "Go is awesome",
			},
			expectedID: 1,
		},
		{
			name:        "repository returns error",
			repoErr:     errors.New("database error"),
			expectedErr: true,
		},
		{
			name:     "returns nil when note does not exist",
			repoNote: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo := &mocks.FakeNoteRepository{
				Note: tt.repoNote,
				Err:  tt.repoErr,
			}

			service := NewNoteService(repo)

			note, err := service.GetByID(tt.expectedID)

			if tt.expectedErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tt.repoNote == nil {
				assert.Nil(t, note)
				return
			}

			assert.NotNil(t, note)
			assert.Equal(t, tt.repoNote.ID, note.ID)
			assert.Equal(t, tt.repoNote.Title, note.Title)
			assert.Equal(t, tt.repoNote.Content, note.Content)
		})
	}
}

func TestNoteServiceCreate(t *testing.T) {

	tests := []struct {
		name        string
		note        *models.Note
		repoErr     error
		expectedErr bool
	}{
		{
			name: "creates note successfully",
			note: &models.Note{
				Title:   "Go",
				Content: "Go is awesome",
			},
		},
		{
			name: "repository returns error",
			note: &models.Note{
				Title:   "Docker",
				Content: "Docker containers",
			},
			repoErr:     errors.New("database error"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo := &mocks.FakeNoteRepository{
				Err: tt.repoErr,
			}

			service := NewNoteService(repo)

			err := service.Create(tt.note)

			if tt.expectedErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tt.note, repo.CreatedNote)
		})
	}
}

func TestNoteServiceUpdate(t *testing.T) {

	tests := []struct {
		name        string
		note        *models.Note
		repoErr     error
		expectedErr bool
	}{
		{
			name: "updates note successfully",
			note: &models.Note{
				ID:      1,
				Title:   "Go Updated",
				Content: "Go is awesome and fast",
			},
		},
		{
			name: "repository returns error",
			note: &models.Note{
				ID:      1,
				Title:   "Go",
				Content: "Go is awesome",
			},
			repoErr:     errors.New("database error"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo := &mocks.FakeNoteRepository{
				Err: tt.repoErr,
			}

			service := NewNoteService(repo)

			err := service.Update(tt.note)

			if tt.expectedErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tt.note, repo.UpdatedNote)
		})
	}
}

func TestNoteServiceDelete(t *testing.T) {

	tests := []struct {
		name        string
		id          int
		repoErr     error
		expectedErr bool
	}{
		{
			name: "deletes note successfully",
			id:   1,
		},
		{
			name:        "repository returns error",
			id:          1,
			repoErr:     errors.New("database error"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo := &mocks.FakeNoteRepository{
				Err: tt.repoErr,
			}

			service := NewNoteService(repo)

			err := service.Delete(tt.id)

			if tt.expectedErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tt.id, repo.DeletedID)
		})
	}
}
