package mocks

type FakeChatService struct {
	Answer string
	Err    error

	LastMessage string
	Called      bool
}

func (f *FakeChatService) Ask(message string) (string, error) {
	f.Called = true
	f.LastMessage = message

	return f.Answer, f.Err
}
