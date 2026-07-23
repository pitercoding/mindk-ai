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
		name               string
		message            string
		context            *models.ChatContext
		notes              []models.Note
		noteErr            error
		llmAnswer          string
		llmErr             error
		historyErr         error
		expectedAnswer     string
		expectError        bool
		expectedQuestion   string
		expectedSavedReply string
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
			llmAnswer:          "Go is awesome",
			expectedAnswer:     "Go is awesome",
			expectedQuestion:   "What do my notes say about Go?",
			expectedSavedReply: "Go is awesome",
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
			name:        "history service returns error",
			message:     "Go",
			notes:       []models.Note{{Title: "Go", Content: "Go is awesome"}},
			llmAnswer:   "Go is awesome",
			historyErr:  errors.New("history error"),
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
			llmAnswer:          "Go is a compiled language",
			expectedAnswer:     "Go is a compiled language",
			expectedQuestion:   "Explain this note",
			expectedSavedReply: "Go is a compiled language",
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

			historyService := &mocks.FakeChatHistoryService{
				Recent: []models.ChatHistory{
					{
						Question: "Previous question",
						Answer:   "Previous answer",
					},
				},
				Err: tt.historyErr,
			}

			chatMessageService := &mocks.FakeChatMessageService{}

			service := NewChatService(
				noteProvider,
				historyService,
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

			assert.Equal(
				t,
				tt.expectedQuestion,
				historyService.SavedQuestion,
			)

			assert.Equal(
				t,
				tt.expectedSavedReply,
				historyService.SavedAnswer,
			)

			assert.Equal(
				t,
				5,
				historyService.LastLimit,
			)

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
