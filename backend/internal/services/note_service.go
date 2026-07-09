package services

import (
	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type NoteRepository interface {
	Create(note *models.Note) error
	GetAll() ([]models.Note, error)
	GetByID(id int) (*models.Note, error)
	Update(note *models.Note) error
	Delete(id int) error
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

func (s *NoteService) GetAll() ([]models.Note, error) {
	return s.repo.GetAll()
}

func (s *NoteService) GetByID(id int) (*models.Note, error) {
	return s.repo.GetByID(id)
}

func (s *NoteService) Update(note *models.Note) error {
	return s.repo.Update(note)
}

func (s *NoteService) Delete(id int) error {
	return s.repo.Delete(id)
}
