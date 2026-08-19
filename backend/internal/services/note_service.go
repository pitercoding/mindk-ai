package services

import (
	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type NoteRepository interface {
	Create(note *models.Note) error
	GetAll(userID string) ([]models.Note, error)
	GetByID(id int, userID string) (*models.Note, error)
	Update(note *models.Note) error
	Delete(id int, userID string) error
}

type NoteService struct {
	repo NoteRepository
}

func NewNoteService(repo NoteRepository) *NoteService {
	return &NoteService{
		repo: repo,
	}
}

func (s *NoteService) Create(note *models.Note) error {
	return s.repo.Create(note)
}

func (s *NoteService) GetAll(userID string) ([]models.Note, error) {
	return s.repo.GetAll(userID)
}

func (s *NoteService) GetByID(id int, userID string) (*models.Note, error) {
	return s.repo.GetByID(id, userID)
}

func (s *NoteService) Update(note *models.Note) error {
	return s.repo.Update(note)
}

func (s *NoteService) Delete(id int, userID string) error {
	return s.repo.Delete(id, userID)
}
