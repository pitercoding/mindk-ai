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
		name string

		// resolveSession inputs
		sessionID   int
		mockSession *models.ChatSession
		sessionErr  error
		mode        string
		noteIDInput *int
		title       string

		message string

		notes   []models.Note
		note    *models.Note
		noteErr error

		messages []models.ChatMessage

		llmAnswer string
		llmErr    error

		documentContext string
		documentSources []models.ChatSource
		documentErr     error

		expectedAnswer    string
		expectedSessionID int
		expectedSources   []models.ChatSource
		expectCreated     bool
		expectError       bool
	}{
		{
			name:      "existing session, knowledge mode returns answer",
			sessionID: 1,

			mockSession: &models.ChatSession{
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

			llmAnswer:         "Go is awesome",
			expectedAnswer:    "Go is awesome",
			expectedSessionID: 1,
		},

		{
			name:      "existing session, note mode returns answer",
			sessionID: 1,

			mockSession: &models.ChatSession{
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

			llmAnswer:         "Go is a compiled language",
			expectedAnswer:    "Go is a compiled language",
			expectedSessionID: 1,
		},

		{
			name:      "session provider returns error",
			sessionID: 1,

			sessionErr: errors.New("database error"),

			message:     "Hello",
			expectError: true,
		},

		{
			name:      "session does not exist",
			sessionID: 1,

			mockSession: nil,

			message:     "Hello",
			expectError: true,
		},

		{
			name:      "note session without note_id returns error",
			sessionID: 1,

			mockSession: &models.ChatSession{
				ID:   1,
				Mode: "note",
			},

			message:     "Explain this",
			expectError: true,
		},

		{
			name:      "note provider returns error",
			sessionID: 1,

			mockSession: &models.ChatSession{
				ID:     1,
				Mode:   "note",
				NoteID: intPtr(10),
			},

			message: "Explain this",

			noteErr:     errors.New("note database error"),
			expectError: true,
		},

		{
			name:      "llm returns error",
			sessionID: 1,

			mockSession: &models.ChatSession{
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
			name:      "uses document context in knowledge mode",
			sessionID: 1,

			mockSession: &models.ChatSession{
				ID:   1,
				Mode: "knowledge",
			},

			message: "What is RAG?",

			documentContext: "RAG combines retrieval and generation.",
			documentSources: []models.ChatSource{
				{DocumentID: 1, Name: "rag.md", Score: 0.9},
			},

			llmAnswer:         "RAG is a technique.",
			expectedAnswer:    "RAG is a technique.",
			expectedSessionID: 1,
			expectedSources: []models.ChatSource{
				{DocumentID: 1, Name: "rag.md", Score: 0.9},
			},
		},

		{
			name:      "document context returns error",
			sessionID: 1,

			mockSession: &models.ChatSession{
				ID:   1,
				Mode: "knowledge",
			},

			message: "What is RAG?",

			documentErr: errors.New("search failed"),

			expectError: true,
		},

		{
			name:      "invalid session mode returns error",
			sessionID: 1,

			mockSession: &models.ChatSession{
				ID:   1,
				Mode: "invalid",
			},

			message:     "Hello",
			expectError: true,
		},

		{
			name:      "auto-creates a session when session_id is not provided",
			sessionID: 0,
			mode:      "knowledge",
			title:     "My chat",

			mockSession: &models.ChatSession{
				ID: 42,
			},

			message: "What do my notes say about Go?",

			notes: []models.Note{
				{
					Title:   "Go",
					Content: "Go is awesome",
				},
			},

			llmAnswer:         "Go is awesome",
			expectedAnswer:    "Go is awesome",
			expectedSessionID: 42,
			expectCreated:     true,
		},

		{
			name:        "auto-create note session forwards note_id",
			sessionID:   0,
			mode:        "note",
			noteIDInput: intPtr(10),

			mockSession: &models.ChatSession{
				ID: 43,
			},

			message: "Explain this note",

			note: &models.Note{
				ID:      10,
				Title:   "Go",
				Content: "Go is a compiled language",
			},

			llmAnswer:         "Go is a compiled language",
			expectedAnswer:    "Go is a compiled language",
			expectedSessionID: 43,
			expectCreated:     true,
		},

		{
			name:      "auto-create without mode returns error",
			sessionID: 0,
			mode:      "",

			message:     "Hello",
			expectError: true,
		},

		{
			name:      "auto-create session creation fails",
			sessionID: 0,
			mode:      "knowledge",

			sessionErr: errors.New("insert failed"),

			message:       "Hello",
			expectError:   true,
			expectCreated: true,
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
				Session: tt.mockSession,
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
				Sources: tt.documentSources,
				Err:     tt.documentErr,
			}

			service := NewChatService(
				noteProvider,
				chatSessionService,
				documentContextProvider,
				chatMessageService,
				llmClient,
			)

			answer, sessionID, sources, err := service.Ask(
				tt.sessionID,
				tt.message,
				tt.mode,
				tt.noteIDInput,
				tt.title,
			)

			if tt.expectCreated {
				require.Len(t, chatSessionService.Created, 1)
			} else {
				assert.Empty(t, chatSessionService.Created)
			}

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

			assert.Equal(
				t,
				tt.expectedSessionID,
				sessionID,
			)

			assert.Equal(
				t,
				tt.expectedSources,
				sources,
			)

			require.Len(
				t,
				chatMessageService.Saved,
				2,
			)

			assert.Equal(
				t,
				tt.expectedSessionID,
				chatMessageService.Saved[0].SessionID,
			)

			assert.Equal(
				t,
				tt.expectedSessionID,
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
