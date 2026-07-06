package services

import (
	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
)

type NoteService struct {
	repo *repository.NoteRepository
}

func NewNoteService(repo *repository.NoteRepository) *NoteService {
	return &NoteService{repo: repo}
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
