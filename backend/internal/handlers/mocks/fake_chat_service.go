package mocks

type FakeChatService struct {
	Answer          string
	ResponseSession int
	Err             error

	LastSessionID int
	LastMessage   string
	LastMode      string
	LastNoteID    *int
	LastTitle     string

	Called bool
}

func (f *FakeChatService) Ask(
	sessionID int,
	message string,
	mode string,
	noteID *int,
	title string,
) (string, int, error) {

	f.Called = true

	f.LastSessionID = sessionID
	f.LastMessage = message
	f.LastMode = mode
	f.LastNoteID = noteID
	f.LastTitle = title

	return f.Answer, f.ResponseSession, f.Err

}
