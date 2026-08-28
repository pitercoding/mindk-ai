package services

import (
	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/utils"
)

type DocumentRepository interface {
	Create(document *models.Document) error
	GetAll(userID string) ([]models.Document, error)
	GetByID(id int, userID string) (*models.Document, error)
	Delete(id int, userID string) error
	Search(query string, userID string) ([]models.Document, error)
}

type ChunkService interface {
	CreateMany(chunks []models.DocumentChunk) error
}

type EmbeddingCreator interface {
	GenerateForChunks(chunks []models.DocumentChunk) error
}

type DocumentService struct {
	repo             DocumentRepository
	chunkService     ChunkService
	embeddingService EmbeddingCreator
}

func NewDocumentService(
	repo DocumentRepository,
	chunkService ChunkService,
	embeddingService EmbeddingCreator,
) *DocumentService {

	return &DocumentService{
		repo:             repo,
		chunkService:     chunkService,
		embeddingService: embeddingService,
	}
}

func (s *DocumentService) Create(
	document *models.Document,
) error {

	err := s.repo.Create(document)

	if err != nil {
		return err
	}

	chunks := utils.SplitIntoChunks(
		document.Content,
		500,
	)

	documentChunks := make(
		[]models.DocumentChunk,
		0,
		len(chunks),
	)

	for index, content := range chunks {

		documentChunks = append(
			documentChunks,
			models.DocumentChunk{
				DocumentID: document.ID,
				ChunkIndex: index,
				Content:    content,
			},
		)
	}

	err = s.chunkService.CreateMany(documentChunks)

	if err != nil {

		// rollback document creation
		s.repo.Delete(document.ID, document.UserID)

		return err
	}

	err = s.embeddingService.GenerateForChunks(documentChunks)

	if err != nil {

		// rollback document creation
		// cascade deletes chunks and embeddings
		s.repo.Delete(document.ID, document.UserID)

		return err
	}

	return nil
}

func (s *DocumentService) GetAll(userID string) ([]models.Document, error) {
	return s.repo.GetAll(userID)
}

func (s *DocumentService) GetByID(
	id int,
	userID string,
) (*models.Document, error) {

	return s.repo.GetByID(id, userID)
}

func (s *DocumentService) Delete(
	id int,
	userID string,
) error {

	return s.repo.Delete(id, userID)
}

func (s *DocumentService) Search(
	query string,
	userID string,
) ([]models.Document, error) {

	return s.repo.Search(query, userID)
}
