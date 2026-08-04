package services

import "fmt"

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

func (s *DocumentContextService) BuildContext(
	query string,
	limit int,
) (string, error) {

	results, err := s.searchService.Search(
		query,
		limit,
	)

	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", nil
	}

	context := "Relevant information from documents:\n\n"

	for _, result := range results {

		context += fmt.Sprintf(
			"Document chunk %d:\n%s\n\n",
			result.ChunkIndex,
			result.Content,
		)
	}

	return context, nil
}
