package services

import (
	"sort"

	"github.com/pitercoding/mindk-ai/backend/internal/utils"
)

type SearchResult struct {
	DocumentID int
	ChunkIndex int
	Content    string
	Score      float32
}

type DocumentSearchService struct {
	embeddings DocumentEmbeddingRepository
	generator  EmbeddingGenerator
}

func NewDocumentSearchService(
	embeddings DocumentEmbeddingRepository,
	generator EmbeddingGenerator,
) *DocumentSearchService {

	return &DocumentSearchService{
		embeddings: embeddings,
		generator:  generator,
	}
}

func (s *DocumentSearchService) Search(
	query string,
	limit int,
) ([]SearchResult, error) {

	queryVector, err := s.generator.Generate(query)
	if err != nil {
		return nil, err
	}

	items, err := s.embeddings.GetAll()
	if err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0)

	for _, item := range items {

		vector, err := utils.ParseVector(item.Embedding)
		if err != nil {
			continue
		}

		score := utils.CosineSimilarity(
			queryVector,
			vector,
		)

		results = append(
			results,
			SearchResult{
				DocumentID: item.DocumentID,
				ChunkIndex: item.ChunkIndex,
				Content:    item.Content,
				Score:      score,
			},
		)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if limit > len(results) {
		limit = len(results)
	}

	return results[:limit], nil
}
