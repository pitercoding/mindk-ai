package mocks

type FakeDocumentContextProvider struct {
	Context string
	Err     error
	Called  bool

	Query string
	Limit int
}

func (f *FakeDocumentContextProvider) BuildContext(
	query string,
	limit int,
) (string, error) {

	f.Called = true

	f.Query = query
	f.Limit = limit

	return f.Context, f.Err
}
