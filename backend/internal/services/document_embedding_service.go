package services

import (
	"encoding/json"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type DocumentEmbeddingRepository interface {
	Create(documentEmbedding *models.DocumentEmbedding) error
	CreateMany([]models.DocumentEmbedding) error
	GetAll(userID string) ([]models.DocumentEmbedding, error)
}

type DocumentEmbeddingService struct {
	repo      DocumentEmbeddingRepository
	generator EmbeddingGenerator
}

func NewDocumentEmbeddingService(
	repo DocumentEmbeddingRepository,
	generator EmbeddingGenerator,
) *DocumentEmbeddingService {

	return &DocumentEmbeddingService{
		repo:      repo,
		generator: generator,
	}
}

func (s *DocumentEmbeddingService) Create(
	documentEmbedding *models.DocumentEmbedding,
) error {

	return s.repo.Create(documentEmbedding)
}

func (s *DocumentEmbeddingService) CreateMany(
	embeddings []models.DocumentEmbedding,
) error {

	return s.repo.CreateMany(embeddings)
}

func (s *DocumentEmbeddingService) GenerateForChunks(
	chunks []models.DocumentChunk,
) error {

	embeddings := make(
		[]models.DocumentEmbedding,
		0,
		len(chunks),
	)

	for _, chunk := range chunks {

		vector, err := s.generator.Generate(
			chunk.Content,
		)

		if err != nil {
			return err
		}

		jsonVector, err := json.Marshal(vector)

		if err != nil {
			return err
		}

		embeddings = append(
			embeddings,
			models.DocumentEmbedding{
				ChunkID:   chunk.ID,
				Embedding: string(jsonVector),
			},
		)
	}

	return s.repo.CreateMany(embeddings)
}

func (s *DocumentEmbeddingService) GetAll(userID string) ([]models.DocumentEmbedding, error) {
	return s.repo.GetAll(userID)
}
