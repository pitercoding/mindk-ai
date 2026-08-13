package services

import (
	"sort"

	"github.com/pitercoding/mindk-ai/backend/internal/utils"
)

// DefaultMinRelevanceScore filters out chunks whose cosine similarity to the query falls below this threshold. Calibrated live against OpenAI's text-embedding-3-small: unrelated content scored ~0.18-0.24 (the noise floor for short documents on this model), while genuinely relevant chunks scored ~0.63-0.69. 0.3 sits well above the noise floor with margin to spare - a bar like 0.7 would discard real matches entirely.
const DefaultMinRelevanceScore = 0.3

type SearchResult struct {
	DocumentID   int
	DocumentName string
	ChunkIndex   int
	Content      string
	Score        float32
}

type DocumentSearchService struct {
	embeddings DocumentEmbeddingRepository
	generator  EmbeddingGenerator
	minScore   float32
}

func NewDocumentSearchService(
	embeddings DocumentEmbeddingRepository,
	generator EmbeddingGenerator,
	minScore float32,
) *DocumentSearchService {

	return &DocumentSearchService{
		embeddings: embeddings,
		generator:  generator,
		minScore:   minScore,
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

		if score < s.minScore {
			continue
		}

		results = append(
			results,
			SearchResult{
				DocumentID:   item.DocumentID,
				DocumentName: item.DocumentName,
				ChunkIndex:   item.ChunkIndex,
				Content:      item.Content,
				Score:        score,
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
