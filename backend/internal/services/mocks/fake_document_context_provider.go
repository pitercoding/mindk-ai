package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeDocumentContextProvider struct {
	Context string
	Sources []models.ChatSource
	Err     error
	Called  bool

	Query  string
	Limit  int
	UserID string
}

func (f *FakeDocumentContextProvider) BuildContext(
	query string,
	limit int,
	userID string,
) (string, []models.ChatSource, error) {

	f.Called = true

	f.Query = query
	f.Limit = limit
	f.UserID = userID

	return f.Context, f.Sources, f.Err
}
