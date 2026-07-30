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

type DocumentService struct {
	repo         DocumentRepository
	chunkService ChunkCreator
}

func NewDocumentService(
	repo DocumentRepository,
	chunkService ChunkCreator,
) *DocumentService {

	return &DocumentService{
		repo:         repo,
		chunkService: chunkService,
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
		5,
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

	return s.chunkService.CreateMany(documentChunks)
}
