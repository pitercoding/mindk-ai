package mocks

type FakeChatService struct {
	Answer string
	Err    error
}

func (f *FakeChatService) Ask(message string) (string, error) {
	return f.Answer, f.Err
}
