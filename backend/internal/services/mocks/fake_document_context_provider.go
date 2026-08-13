package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeDocumentContextProvider struct {
	Context string
	Sources []models.ChatSource
	Err     error
	Called  bool

	Query string
	Limit int
}

func (f *FakeDocumentContextProvider) BuildContext(
	query string,
	limit int,
) (string, []models.ChatSource, error) {

	f.Called = true

	f.Query = query
	f.Limit = limit

	return f.Context, f.Sources, f.Err
}
