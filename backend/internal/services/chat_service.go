package services

import (
	"github.com/pitercoding/mindk-ai/backend/internal/llm"
)

type ChatService struct {
	noteService *NoteService
}

func NewChatService(noteService *NoteService) *ChatService {
	return &ChatService{
		noteService: noteService,
	}
}

func (s *ChatService) Ask(message string) (string, error) {
	notes, err := s.noteService.GetAll()
	if err != nil {
		return "", err
	}

	prompt := llm.BuildPrompt(message, notes)

	return prompt, nil
}
