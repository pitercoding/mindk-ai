package services

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type DocumentChunkRepository interface {
	Create(documentChunk *models.DocumentChunk) error
	CreateMany([]models.DocumentChunk) error
	GetByDocumentID(documentID int) ([]models.DocumentChunk, error)
	DeleteByDocumentID(documentID int) error
}

type DocumentChunkService struct {
	repo DocumentChunkRepository
}

func NewDocumentChunkService(repo DocumentChunkRepository) *DocumentChunkService {
	return &DocumentChunkService{repo: repo}
}

func (s *DocumentChunkService) Create(documentChunk *models.DocumentChunk) error {
	return s.repo.Create(documentChunk)
}

func (s *DocumentChunkService) CreateMany(chunks []models.DocumentChunk) error {
	return s.repo.CreateMany(chunks)
}

func (s *DocumentChunkService) GetByDocumentID(documentID int) ([]models.DocumentChunk, error) {
	return s.repo.GetByDocumentID(documentID)
}

func (s *DocumentChunkService) DeleteByDocumentID(documentID int) error {
	return s.repo.DeleteByDocumentID(documentID)
}
