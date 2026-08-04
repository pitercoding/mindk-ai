package services

import (
	"errors"
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
}

func (f *fakeDocumentSearcher) Search(
	query string,
	limit int,
) ([]SearchResult, error) {

	f.Called = true
	f.Query = query
	f.Limit = limit

	if f.Err != nil {
		return nil, f.Err
	}

	return f.Results, nil
}

func TestDocumentContextServiceBuildContext(t *testing.T) {

	searcher := &fakeDocumentSearcher{
		Results: []SearchResult{
			{
				DocumentID: 1,
				ChunkIndex: 0,
				Content:    "Go is a programming language.",
				Score:      1,
			},
			{
				DocumentID: 2,
				ChunkIndex: 3,
				Content:    "Go uses goroutines for concurrency.",
				Score:      0.8,
			},
		},
	}

	service := NewDocumentContextService(
		searcher,
	)

	context, err := service.BuildContext(
		"How does Go work?",
		2,
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
		"Go is a programming language.",
	)

	assert.Contains(
		t,
		context,
		"Go uses goroutines for concurrency.",
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

func TestDocumentContextServiceBuildContext_NoResults(t *testing.T) {

	searcher := &fakeDocumentSearcher{
		Results: []SearchResult{},
	}

	service := NewDocumentContextService(
		searcher,
	)

	context, err := service.BuildContext(
		"unknown topic",
		5,
	)

	require.NoError(
		t,
		err,
	)

	assert.Empty(
		t,
		context,
	)
}

func TestDocumentContextServiceBuildContext_SearchError(t *testing.T) {

	searcher := &fakeDocumentSearcher{
		Err: errors.New("search failed"),
	}

	service := NewDocumentContextService(
		searcher,
	)

	context, err := service.BuildContext(
		"test",
		5,
	)

	assert.Error(
		t,
		err,
	)

	assert.Empty(
		t,
		context,
	)
}
