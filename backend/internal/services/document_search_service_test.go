package services

import (
	"errors"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocumentSearchServiceSearch(t *testing.T) {

	repository := &mocks.FakeDocumentEmbeddingRepository{
		Embeddings: []models.DocumentEmbedding{
			{
				ChunkID:    1,
				DocumentID: 10,
				ChunkIndex: 0,
				Content:    "Go is a programming language",
				Embedding:  `[1,0,0]`,
			},
			{
				ChunkID:    2,
				DocumentID: 20,
				ChunkIndex: 1,
				Content:    "Java is another language",
				Embedding:  `[0,1,0]`,
			},
			{
				ChunkID:    3,
				DocumentID: 30,
				ChunkIndex: 2,
				Content:    "React is a frontend library",
				Embedding:  `[0.9,0.1,0]`,
			},
		},
	}

	generator := &mocks.FakeEmbeddingGenerator{
		Vector: []float32{
			1,
			0,
			0,
		},
	}

	service := NewDocumentSearchService(
		repository,
		generator,
	)

	results, err := service.Search(
		"programming language",
		2,
	)

	require.NoError(
		t,
		err,
	)

	require.Len(
		t,
		results,
		2,
	)

	// First result should be the most similar
	assert.Equal(
		t,
		10,
		results[0].DocumentID,
	)

	assert.Equal(
		t,
		0,
		results[0].ChunkIndex,
	)

	// Second result should be React because
	// vector [0.9,0.1,0] is closer than [0,1,0]
	assert.Equal(
		t,
		30,
		results[1].DocumentID,
	)

	assert.True(
		t,
		generator.Called,
	)

	assert.Equal(
		t,
		"programming language",
		generator.Texts[0],
	)
}

func TestDocumentSearchServiceSearch_LimitGreaterThanResults(t *testing.T) {

	repository := &mocks.FakeDocumentEmbeddingRepository{
		Embeddings: []models.DocumentEmbedding{
			{
				DocumentID: 1,
				Content:    "test",
				Embedding:  `[1,0]`,
			},
		},
	}

	generator := &mocks.FakeEmbeddingGenerator{
		Vector: []float32{
			1,
			0,
		},
	}

	service := NewDocumentSearchService(
		repository,
		generator,
	)

	results, err := service.Search(
		"test",
		10,
	)

	require.NoError(
		t,
		err,
	)

	assert.Len(
		t,
		results,
		1,
	)
}

func TestDocumentSearchServiceSearch_EmbeddingGeneratorError(t *testing.T) {

	repository := &mocks.FakeDocumentEmbeddingRepository{}

	generator := &mocks.FakeEmbeddingGenerator{
		Err: errors.New("embedding failed"),
	}

	service := NewDocumentSearchService(
		repository,
		generator,
	)

	results, err := service.Search(
		"test",
		5,
	)

	assert.Error(
		t,
		err,
	)

	assert.Nil(
		t,
		results,
	)
}

func TestDocumentSearchServiceSearch_RepositoryError(t *testing.T) {

	repository := &mocks.FakeDocumentEmbeddingRepository{
		Err: errors.New("database error"),
	}

	generator := &mocks.FakeEmbeddingGenerator{
		Vector: []float32{
			1,
			0,
		},
	}

	service := NewDocumentSearchService(
		repository,
		generator,
	)

	results, err := service.Search(
		"test",
		5,
	)

	assert.Error(
		t,
		err,
	)

	assert.Nil(
		t,
		results,
	)
}
