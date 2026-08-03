package services

import (
	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/utils"
)

type DocumentRepository interface {
	Create(document *models.Document) error
	GetAll() ([]models.Document, error)
	GetByID(id int) (*models.Document, error)
	Delete(id int) error
	Search(query string) ([]models.Document, error)
}

type ChunkCreator interface {
	CreateMany(chunks []models.DocumentChunk) error
}

type EmbeddingCreator interface {
	GenerateForChunks(chunks []models.DocumentChunk) error
}

type DocumentService struct {
	repo             DocumentRepository
	chunkService     ChunkCreator
	embeddingService EmbeddingCreator
}

func NewDocumentService(
	repo DocumentRepository,
	chunkService ChunkCreator,
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
		return err
	}

	return s.embeddingService.GenerateForChunks(documentChunks)
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
