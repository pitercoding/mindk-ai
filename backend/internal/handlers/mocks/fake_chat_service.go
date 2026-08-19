package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeChatService struct {
	Answer          string
	ResponseSession int
	Sources         []models.ChatSource
	Err             error

	LastUserID    string
	LastSessionID int
	LastMessage   string
	LastMode      string
	LastNoteID    *int
	LastTitle     string

	Called bool
}

func (f *FakeChatService) Ask(
	userID string,
	sessionID int,
	message string,
	mode string,
	noteID *int,
	title string,
) (string, int, []models.ChatSource, error) {

	f.Called = true

	f.LastUserID = userID
	f.LastSessionID = sessionID
	f.LastMessage = message
	f.LastMode = mode
	f.LastNoteID = noteID
	f.LastTitle = title

	return f.Answer, f.ResponseSession, f.Sources, f.Err

}
