package services

import (
	"errors"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatHistoryServiceSave(t *testing.T) {

	repo := &mocks.FakeChatHistoryRepository{}

	service := NewChatHistoryService(repo)

	err := service.Save(
		"What is Docker?",
		"Docker is a container platform",
	)

	require.NoError(t, err)

	require.NotNil(
		t,
		repo.CreatedHistory,
	)

	assert.Equal(
		t,
		"What is Docker?",
		repo.CreatedHistory.Question,
	)

	assert.Equal(
		t,
		"Docker is a container platform",
		repo.CreatedHistory.Answer,
	)
}

func TestChatHistoryServiceSave_Error(t *testing.T) {

	repo := &mocks.FakeChatHistoryRepository{
		Err: errors.New("database error"),
	}

	service := NewChatHistoryService(repo)

	err := service.Save(
		"What is Docker?",
		"Docker is a container platform",
	)

	assert.Error(t, err)

	require.NotNil(
		t,
		repo.CreatedHistory,
	)

	assert.Equal(
		t,
		"What is Docker?",
		repo.CreatedHistory.Question,
	)
}

func TestChatHistoryServiceGetAll(t *testing.T) {

	repo := &mocks.FakeChatHistoryRepository{
		History: []models.ChatHistory{
			{
				ID:       1,
				Question: "What is Go?",
				Answer:   "Go is a programming language",
			},
			{
				ID:       2,
				Question: "What is Docker?",
				Answer:   "Docker is a container platform",
			},
		},
		Total: 10,
	}

	service := NewChatHistoryService(repo)

	history, total, err := service.GetAll(
		2,
		5,
	)

	require.NoError(t, err)

	assert.Len(
		t,
		history,
		2,
	)

	assert.Equal(
		t,
		10,
		total,
	)

	assert.Equal(
		t,
		5,
		repo.LastLimit,
	)

	assert.Equal(
		t,
		5,
		repo.LastOffset,
	)
}

func TestChatHistoryServiceGetAll_Error(t *testing.T) {

	repo := &mocks.FakeChatHistoryRepository{
		Err: errors.New("database error"),
	}

	service := NewChatHistoryService(repo)

	history, total, err := service.GetAll(
		1,
		10,
	)

	assert.Error(t, err)

	assert.Nil(
		t,
		history,
	)

	assert.Equal(
		t,
		0,
		total,
	)
}

func TestChatHistoryServiceGetRecent(t *testing.T) {

	repo := &mocks.FakeChatHistoryRepository{
		Recent: []models.ChatHistory{
			{
				ID:       2,
				Question: "What is Docker?",
				Answer:   "Docker is a container platform",
			},
			{
				ID:       1,
				Question: "What is Go?",
				Answer:   "Go is a programming language",
			},
		},
	}

	service := NewChatHistoryService(repo)

	history, err := service.GetRecent(5)

	require.NoError(t, err)

	assert.Len(
		t,
		history,
		2,
	)

	assert.Equal(
		t,
		5,
		repo.LastRecentLimit,
	)

	assert.Equal(
		t,
		"What is Docker?",
		history[0].Question,
	)

	assert.Equal(
		t,
		"What is Go?",
		history[1].Question,
	)
}

func TestChatHistoryServiceGetRecent_Error(t *testing.T) {

	repo := &mocks.FakeChatHistoryRepository{
		Err: errors.New("database error"),
	}

	service := NewChatHistoryService(repo)

	history, err := service.GetRecent(5)

	assert.Error(t, err)

	assert.Nil(
		t,
		history,
	)
}
