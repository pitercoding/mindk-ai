package services

import (
	"fmt"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

// maxContextChars caps the total size of the assembled document context. A character count is used instead of a token count to avoid pulling in a tokenizer dependency; it's a conservative proxy (tokens are ~4 chars on average, so this budgets for roughly 2k-4k tokens of context).
const maxContextChars = 8000

type DocumentSearcher interface {
	Search(query string, limit int) ([]SearchResult, error)
}

type DocumentContextService struct {
	searchService DocumentSearcher
}

func NewDocumentContextService(searchService DocumentSearcher) *DocumentContextService {

	return &DocumentContextService{
		searchService: searchService,
	}
}

// BuildContext retrieves the chunks most relevant to query and assembles them into prompt-ready text, stopping once maxContextChars is reached. It also returns the distinct documents actually included, highest score first, so callers can surface them as citations.
func (s *DocumentContextService) BuildContext(
	query string,
	limit int,
) (string, []models.ChatSource, error) {

	results, err := s.searchService.Search(
		query,
		limit,
	)

	if err != nil {
		return "", nil, err
	}

	if len(results) == 0 {
		return "", nil, nil
	}

	context := "Relevant information from documents:\n\n"
	contextSize := len(context)

	sources := make([]models.ChatSource, 0, len(results))
	seenDocuments := make(map[int]bool)
	chunksAdded := 0

	for _, result := range results {

		chunk := fmt.Sprintf(
			"[Document: %s]\nChunk %d:\n%s\n\n",
			result.DocumentName,
			result.ChunkIndex,
			result.Content,
		)

		// The single most relevant chunk is always included, even if it alone exceeds the budget, so a legitimately relevant answer is never dropped entirely; only chunks beyond it are budget-bound.
		if chunksAdded > 0 && contextSize+len(chunk) > maxContextChars {
			break
		}

		context += chunk
		contextSize += len(chunk)
		chunksAdded++

		if !seenDocuments[result.DocumentID] {

			seenDocuments[result.DocumentID] = true

			sources = append(
				sources,
				models.ChatSource{
					DocumentID: result.DocumentID,
					Name:       result.DocumentName,
					Score:      result.Score,
				},
			)
		}
	}

	return context, sources, nil
}
