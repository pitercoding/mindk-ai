package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeChatService struct {
	Answer string
	Err    error

	LastMessage string
	LastContext *models.ChatContext

	Called bool
}

func (f *FakeChatService) Ask(
	message string,
	context *models.ChatContext,
) (string, error) {

	f.Called = true

	f.LastMessage = message
	f.LastContext = context

	return f.Answer, f.Err
}
