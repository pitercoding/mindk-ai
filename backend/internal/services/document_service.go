package services

import (
	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type DocumentRepository interface {
	Create(document *models.Document) error
	GetAll() ([]models.Document, error)
	GetByID(id int) (*models.Document, error)
	Delete(id int) error
	Search(query string) ([]models.Document, error)
}

type DocumentService struct {
	repo DocumentRepository
}

func NewDocumentService(repo DocumentRepository) *DocumentService {
	return &DocumentService{
		repo: repo,
	}
}

func (s *DocumentService) Create(document *models.Document) error {
	return s.repo.Create(document)
}

func (s *DocumentService) GetAll() ([]models.Document, error) {
	return s.repo.GetAll()
}

func (s *DocumentService) GetByID(id int) (*models.Document, error) {
	return s.repo.GetByID(id)
}

func (s *DocumentService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *DocumentService) Search(query string) ([]models.Document, error) {
	return s.repo.Search(query)
}
