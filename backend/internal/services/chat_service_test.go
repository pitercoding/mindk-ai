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
		session         *models.ChatSession
		sessionErr      error
		message         string
		notes           []models.Note
		note            *models.Note
		noteErr         error
		messages        []models.ChatMessage
		llmAnswer       string
		llmErr          error
		documentContext string
		documentErr     error
		expectedAnswer  string
		expectError     bool
	}{
		{
			name: "knowledge session returns answer successfully",

			session: &models.ChatSession{
				ID:   1,
				Mode: "knowledge",
			},

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
			name: "note session returns answer successfully",

			session: &models.ChatSession{
				ID:     1,
				Mode:   "note",
				NoteID: intPtr(10),
			},

			message: "Explain this note",

			note: &models.Note{
				ID:      10,
				Title:   "Go",
				Content: "Go is a compiled language",
			},

			llmAnswer:      "Go is a compiled language",
			expectedAnswer: "Go is a compiled language",
		},

		{
			name: "session provider returns error",

			sessionErr: errors.New("database error"),

			message:     "Hello",
			expectError: true,
		},

		{
			name: "session does not exist",

			session: nil,

			message:     "Hello",
			expectError: true,
		},

		{
			name: "note session without note returns error",

			session: &models.ChatSession{
				ID:   1,
				Mode: "note",
			},

			message:     "Explain this",
			expectError: true,
		},

		{
			name: "note provider returns error",

			session: &models.ChatSession{
				ID:     1,
				Mode:   "note",
				NoteID: intPtr(10),
			},

			message: "Explain this",

			noteErr:     errors.New("note database error"),
			expectError: true,
		},

		{
			name: "llm returns error",

			session: &models.ChatSession{
				ID:   1,
				Mode: "knowledge",
			},

			message: "What is Go?",

			notes: []models.Note{
				{
					Title:   "Go",
					Content: "Go is awesome",
				},
			},

			llmErr:      errors.New("openai error"),
			expectError: true,
		},

		{
			name: "uses document context in knowledge mode",

			session: &models.ChatSession{
				ID:   1,
				Mode: "knowledge",
			},

			message: "What is RAG?",

			documentContext: "RAG combines retrieval and generation.",

			llmAnswer:      "RAG is a technique.",
			expectedAnswer: "RAG is a technique.",
		},

		{
			name: "document context returns error",

			session: &models.ChatSession{
				ID:   1,
				Mode: "knowledge",
			},

			message: "What is RAG?",

			documentErr: errors.New("search failed"),

			expectError: true,
		},

		{
			name: "invalid session mode returns error",

			session: &models.ChatSession{
				ID:   1,
				Mode: "invalid",
			},

			message:     "Hello",
			expectError: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			noteProvider := &mocks.FakeNoteProvider{
				Notes: tt.notes,
				Note:  tt.note,
				Err:   tt.noteErr,
			}

			chatSessionService := &mocks.FakeChatSessionService{
				Session: tt.session,
				Err:     tt.sessionErr,
			}

			chatMessageService := &mocks.FakeChatMessageService{
				Messages: tt.messages,
			}

			llmClient := &mocks.FakeLLMClient{
				Response: tt.llmAnswer,
				Err:      tt.llmErr,
			}

			documentContextProvider := &mocks.FakeDocumentContextProvider{
				Context: tt.documentContext,
				Err:     tt.documentErr,
			}

			service := NewChatService(
				noteProvider,
				chatSessionService,
				documentContextProvider,
				chatMessageService,
				llmClient,
			)

			answer, err := service.Ask(
				1,
				tt.message,
			)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(
				t,
				tt.expectedAnswer,
				answer,
			)

			require.Len(
				t,
				chatMessageService.Saved,
				2,
			)

			assert.Equal(
				t,
				1,
				chatMessageService.Saved[0].SessionID,
			)

			assert.Equal(
				t,
				1,
				chatMessageService.Saved[1].SessionID,
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
		})
	}
}

func intPtr(value int) *int {
	return &value
}
