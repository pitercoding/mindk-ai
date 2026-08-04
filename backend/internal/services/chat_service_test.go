package services

import (
	"errors"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatServiceAsk(t *testing.T) {

	tests := []struct {
		name            string
		message         string
		context         *models.ChatContext
		notes           []models.Note
		noteErr         error
		llmAnswer       string
		llmErr          error
		expectedAnswer  string
		expectError     bool
		documentContext string
		documentErr     error
	}{
		{
			name:    "returns answer successfully",
			message: "What do my notes say about Go?",
			notes: []models.Note{
				{
					Title:   "Go",
					Content: "Go is awesome",
				},
			},
			llmAnswer:      "Go is awesome",
			expectedAnswer: "Go is awesome",
		},
		{
			name:        "note provider returns error",
			message:     "Go",
			noteErr:     errors.New("database error"),
			expectError: true,
		},
		{
			name:        "llm returns error",
			message:     "Go",
			notes:       []models.Note{{Title: "Go", Content: "Go is awesome"}},
			llmErr:      errors.New("openai error"),
			expectError: true,
		},
		{
			name:    "saves chat messages when note context exists",
			message: "Explain this note",
			context: &models.ChatContext{
				NoteID:  1,
				Title:   "Go",
				Content: "Go is a compiled language",
			},
			llmAnswer:      "Go is a compiled language",
			expectedAnswer: "Go is a compiled language",
		},
		{
			name:            "uses document context when available",
			message:         "What is RAG?",
			documentContext: "RAG combines retrieval and generation.",
			llmAnswer:       "RAG is a technique.",
			expectedAnswer:  "RAG is a technique.",
		},
		{
			name:        "document context returns error",
			message:     "What is RAG?",
			documentErr: errors.New("search failed"),
			expectError: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			noteProvider := &mocks.FakeNoteProvider{
				Notes: tt.notes,
				Err:   tt.noteErr,
			}

			llmClient := &mocks.FakeLLMClient{
				Response: tt.llmAnswer,
				Err:      tt.llmErr,
			}

			chatMessageService := &mocks.FakeChatMessageService{}

			documentContextProvider := &mocks.FakeDocumentContextProvider{
				Context: tt.documentContext,
				Err:     tt.documentErr,
			}

			service := NewChatService(
				noteProvider,
				documentContextProvider,
				chatMessageService,
				llmClient,
			)

			answer, err := service.Ask(
				tt.message,
				tt.context,
			)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tt.expectedAnswer, answer)

			if tt.documentContext != "" || tt.documentErr != nil {

				assert.True(
					t,
					documentContextProvider.Called,
				)

				assert.Equal(
					t,
					tt.message,
					documentContextProvider.Query,
				)

				assert.Equal(
					t,
					5,
					documentContextProvider.Limit,
				)
			}

			if tt.context != nil {

				require.Len(
					t,
					chatMessageService.Saved,
					2,
				)

				assert.Equal(
					t,
					tt.context.NoteID,
					chatMessageService.Saved[0].NoteID,
				)

				assert.Equal(
					t,
					tt.context.NoteID,
					chatMessageService.Saved[1].NoteID,
				)

				assert.Equal(
					t,
					"user",
					chatMessageService.Saved[0].Role,
				)

				assert.Equal(
					t,
					tt.message,
					chatMessageService.Saved[0].Content,
				)

				assert.Equal(
					t,
					"assistant",
					chatMessageService.Saved[1].Role,
				)

				assert.Equal(
					t,
					tt.expectedAnswer,
					chatMessageService.Saved[1].Content,
				)
			}
		})
	}
}
