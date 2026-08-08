package mocks

type FakeChatService struct {
	Answer string
	Err    error

	LastSessionID int
	LastMessage   string

	Called bool
}

func (f *FakeChatService) Ask(
	sessionID int,
	message string,
) (string, error) {

	f.Called = true

	f.LastSessionID = sessionID
	f.LastMessage = message

	return f.Answer, f.Err

}
