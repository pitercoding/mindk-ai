package services

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeDocumentSearcher struct {
	Results []SearchResult
	Err     error

	Called bool
	Query  string
	Limit  int
	UserID string
}

func (f *fakeDocumentSearcher) Search(
	query string,
	limit int,
	userID string,
) ([]SearchResult, error) {

	f.Called = true
	f.Query = query
	f.Limit = limit
	f.UserID = userID

	if f.Err != nil {
		return nil, f.Err
	}

	return f.Results, nil
}

func TestDocumentContextServiceBuildContext(t *testing.T) {

	searcher := &fakeDocumentSearcher{
		Results: []SearchResult{
			{
				DocumentID:   1,
				DocumentName: "go-guide.md",
				ChunkIndex:   0,
				Content:      "Go is a programming language.",
				Score:        1,
			},
			{
				DocumentID:   2,
				DocumentName: "concurrency.md",
				ChunkIndex:   3,
				Content:      "Go uses goroutines for concurrency.",
				Score:        0.8,
			},
		},
	}

	service := NewDocumentContextService(
		searcher,
	)

	context, sources, err := service.BuildContext(
		"How does Go work?",
		2,
		"user_1",
	)

	require.NoError(
		t,
		err,
	)

	assert.Contains(
		t,
		context,
		"Relevant information from documents:",
	)

	assert.Contains(
		t,
		context,
		"[Document: go-guide.md]",
	)

	assert.Contains(
		t,
		context,
		"Go is a programming language.",
	)

	assert.Contains(
		t,
		context,
		"[Document: concurrency.md]",
	)

	assert.Contains(
		t,
		context,
		"Go uses goroutines for concurrency.",
	)

	require.Len(
		t,
		sources,
		2,
	)

	assert.Equal(
		t,
		1,
		sources[0].DocumentID,
	)

	assert.Equal(
		t,
		"go-guide.md",
		sources[0].Name,
	)

	assert.Equal(
		t,
		float32(1),
		sources[0].Score,
	)

	assert.Equal(
		t,
		2,
		sources[1].DocumentID,
	)

	assert.True(
		t,
		searcher.Called,
	)

	assert.Equal(
		t,
		"How does Go work?",
		searcher.Query,
	)

	assert.Equal(
		t,
		2,
		searcher.Limit,
	)
}

func TestDocumentContextServiceBuildContext_DedupesSourcesByDocument(t *testing.T) {

	searcher := &fakeDocumentSearcher{
		Results: []SearchResult{
			{
				DocumentID:   1,
				DocumentName: "go-guide.md",
				ChunkIndex:   0,
				Content:      "Go is a programming language.",
				Score:        0.9,
			},
			{
				DocumentID:   1,
				DocumentName: "go-guide.md",
				ChunkIndex:   1,
				Content:      "Go uses goroutines for concurrency.",
				Score:        0.7,
			},
		},
	}

	service := NewDocumentContextService(
		searcher,
	)

	_, sources, err := service.BuildContext(
		"How does Go work?",
		2,
		"user_1",
	)

	require.NoError(
		t,
		err,
	)

	require.Len(
		t,
		sources,
		1,
	)

	assert.Equal(
		t,
		1,
		sources[0].DocumentID,
	)

	// keeps the score of the first (highest-ranked) occurrence
	assert.Equal(
		t,
		float32(0.9),
		sources[0].Score,
	)
}

func TestDocumentContextServiceBuildContext_StopsAtCharBudget(t *testing.T) {

	// Large enough on its own to exceed maxContextChars once the header and wrapper are added, but still gets included because it's the first (most relevant) chunk. The second chunk is the one that gets dropped.
	huge := strings.Repeat("a", maxContextChars)

	searcher := &fakeDocumentSearcher{
		Results: []SearchResult{
			{
				DocumentID:   1,
				DocumentName: "big.md",
				ChunkIndex:   0,
				Content:      huge,
				Score:        1,
			},
			{
				DocumentID:   2,
				DocumentName: "small.md",
				ChunkIndex:   0,
				Content:      "should not fit",
				Score:        0.9,
			},
		},
	}

	service := NewDocumentContextService(
		searcher,
	)

	context, sources, err := service.BuildContext(
		"test",
		2,
		"user_1",
	)

	require.NoError(
		t,
		err,
	)

	assert.Contains(
		t,
		context,
		"[Document: big.md]",
	)

	assert.NotContains(
		t,
		context,
		"should not fit",
	)

	require.Len(
		t,
		sources,
		1,
	)

	assert.Equal(
		t,
		1,
		sources[0].DocumentID,
	)
}

func TestDocumentContextServiceBuildContext_NoResults(t *testing.T) {

	searcher := &fakeDocumentSearcher{
		Results: []SearchResult{},
	}

	service := NewDocumentContextService(
		searcher,
	)

	context, sources, err := service.BuildContext(
		"unknown topic",
		5,
		"user_1",
	)

	require.NoError(
		t,
		err,
	)

	assert.Empty(
		t,
		context,
	)

	assert.Empty(
		t,
		sources,
	)
}

func TestDocumentContextServiceBuildContext_SearchError(t *testing.T) {

	searcher := &fakeDocumentSearcher{
		Err: errors.New("search failed"),
	}

	service := NewDocumentContextService(
		searcher,
	)

	context, sources, err := service.BuildContext(
		"test",
		5,
		"user_1",
	)

	assert.Error(
		t,
		err,
	)

	assert.Empty(
		t,
		context,
	)

	assert.Empty(
		t,
		sources,
	)
}
