package mocks

import (
	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type FakeChatHistoryRepository struct {
	History []models.ChatHistory
	Recent  []models.ChatHistory
	Total   int

	Err error

	CreatedHistory *models.ChatHistory
	Deleted        bool

	LastLimit       int
	LastOffset      int
	LastRecentLimit int
}

func (f *FakeChatHistoryRepository) Create(history *models.ChatHistory) error {
	f.CreatedHistory = history
	return f.Err
}

func (f *FakeChatHistoryRepository) GetAll(limit, offset int) ([]models.ChatHistory, error) {

	f.LastLimit = limit
	f.LastOffset = offset

	return f.History, f.Err
}

func (f *FakeChatHistoryRepository) GetRecent(limit int) ([]models.ChatHistory, error) {

	f.LastRecentLimit = limit

	return f.Recent, f.Err
}

func (f *FakeChatHistoryRepository) Count() (int, error) {
	return f.Total, f.Err
}

func (f *FakeChatHistoryRepository) DeleteAll() error {
	f.Deleted = true
	return f.Err
}
