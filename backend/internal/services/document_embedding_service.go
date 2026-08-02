package services

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type DocumentEmbeddingRepository interface {
	Create(documentEmbedding *models.DocumentEmbedding) error
	CreateMany([]models.DocumentEmbedding) error
	GetByChunkID(chunkID int) (*models.DocumentEmbedding, error)
	DeleteByChunkID(chunkID int) error
}

type DocumentEmbeddingService struct {
	repo DocumentEmbeddingRepository
}

func NewDocumentEmbeddingService(repo DocumentEmbeddingRepository) *DocumentEmbeddingService {
	return &DocumentEmbeddingService{repo: repo}
}

func (s *DocumentEmbeddingService) Create(documentEmbedding *models.DocumentEmbedding) error {
	return s.repo.Create(documentEmbedding)
}

func (s *DocumentEmbeddingService) CreateMany(embeddings []models.DocumentEmbedding) error {
	return s.repo.CreateMany(embeddings)
}

func (s *DocumentEmbeddingService) GetByChunkID(chunkID int) (*models.DocumentEmbedding, error) {
	return s.repo.GetByChunkID(chunkID)
}

func (s *DocumentEmbeddingService) DeleteByChunkID(chunkID int) error {
	return s.repo.DeleteByChunkID(chunkID)
}
