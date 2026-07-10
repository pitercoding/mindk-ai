package llm

import (
	"strings"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestBuildPrompt(t *testing.T) {

	longContent := strings.Repeat("Go is awesome. ", 50)

	tests := []struct {
		name             string
		question         string
		notes            []models.Note
		expectedContains []string
	}{
		{
			name:     "empty notes",
			question: "What is Go?",
			notes:    []models.Note{},
			expectedContains: []string{
				"You are an AI assistant.",
				"NOTES:",
				"USER QUESTION:",
				"What is Go?",
			},
		},
		{
			name:     "one note",
			question: "Tell me about Go",
			notes: []models.Note{
				{
					Title:   "Go",
					Content: "Go is fast.",
				},
			},
			expectedContains: []string{
				"Title: Go",
				"Content: Go is fast.",
				"Tell me about Go",
			},
		},
		{
			name:     "multiple notes",
			question: "Programming languages",
			notes: []models.Note{
				{
					Title:   "Go",
					Content: "Compiled language.",
				},
				{
					Title:   "Java",
					Content: "Object-oriented language.",
				},
			},
			expectedContains: []string{
				"Title: Go",
				"Compiled language.",
				"Title: Java",
				"Object-oriented language.",
			},
		},
		{
			name:     "empty question",
			question: "",
			notes: []models.Note{
				{
					Title:   "Docker",
					Content: "Containers.",
				},
			},
			expectedContains: []string{
				"USER QUESTION:",
				"Title: Docker",
			},
		},
		{
			name:     "long content",
			question: "Explain",
			notes: []models.Note{
				{
					Title:   "Long",
					Content: longContent,
				},
			},
			expectedContains: []string{
				longContent,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			prompt := BuildPrompt(tt.question, tt.notes)

			for _, expected := range tt.expectedContains {
				assert.Contains(t, prompt, expected)
			}
		})
	}
}