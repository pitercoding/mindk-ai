package llm

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

func TestContextBuilder_Build_FullContext(t *testing.T) {
	builder := NewContextBuilder()

	notes := []models.Note{
		{Title: "Project Roadmap", Content: "Ship the RAG pipeline by Q3."},
	}

	messages := []models.ChatMessage{
		{Role: "user", Content: "What is the roadmap?"},
		{Role: "assistant", Content: "The roadmap targets Q3."},
	}

	documentContext := "Source: spec.pdf\nThe API uses REST."

	prompt := builder.Build(
		"When is the deadline?",
		notes,
		messages,
		documentContext,
	)

	assert.Contains(t, prompt, "SYSTEM INSTRUCTIONS:")
	assert.Contains(t, prompt, "CONVERSATION HISTORY:")
	assert.Contains(t, prompt, "KNOWLEDGE CONTEXT:")
	assert.Contains(t, prompt, "DOCUMENT CONTEXT:")
	assert.Contains(t, prompt, "USER QUESTION:")
	assert.Contains(t, prompt, "RESPONSE:")

	assert.Contains(t, prompt, "Title: Project Roadmap")
	assert.Contains(t, prompt, "Ship the RAG pipeline by Q3.")

	assert.Contains(t, prompt, "User:\nWhat is the roadmap?")
	assert.Contains(t, prompt, "Assistant:\nThe roadmap targets Q3.")

	assert.Contains(t, prompt, documentContext)

	assert.Contains(t, prompt, "When is the deadline?")

	assert.True(t,
		strings.Index(prompt, "CONVERSATION HISTORY:") < strings.Index(prompt, "KNOWLEDGE CONTEXT:"),
		"conversation history should come before knowledge context",
	)
	assert.True(t,
		strings.Index(prompt, "KNOWLEDGE CONTEXT:") < strings.Index(prompt, "DOCUMENT CONTEXT:"),
		"knowledge context should come before document context",
	)
	assert.True(t,
		strings.Index(prompt, "DOCUMENT CONTEXT:") < strings.Index(prompt, "USER QUESTION:"),
		"document context should come before the user question",
	)
}

func TestContextBuilder_Build_EmptyContext(t *testing.T) {
	builder := NewContextBuilder()

	prompt := builder.Build(
		"Anything in the knowledge base?",
		[]models.Note{},
		[]models.ChatMessage{},
		"",
	)

	assert.Contains(t, prompt, "No previous conversation.")
	assert.Contains(t, prompt, "No notes available.")
	assert.Contains(t, prompt, "No documents available.")
	assert.Contains(t, prompt, "Anything in the knowledge base?")

	assert.NotContains(t, prompt, "Title:")
	assert.NotContains(t, prompt, "User:")
	assert.NotContains(t, prompt, "Assistant:")
}

func TestContextBuilder_Build_PartialContext(t *testing.T) {
	builder := NewContextBuilder()

	notes := []models.Note{
		{Title: "Meeting Notes", Content: "Discussed Q3 priorities."},
	}

	documentContext := "Source: handbook.pdf\nRemote work policy applies."

	prompt := builder.Build(
		"What did we discuss?",
		notes,
		[]models.ChatMessage{},
		documentContext,
	)

	assert.Contains(t, prompt, "No previous conversation.")

	assert.Contains(t, prompt, "Title: Meeting Notes")
	assert.Contains(t, prompt, "Discussed Q3 priorities.")
	assert.NotContains(t, prompt, "No notes available.")

	assert.Contains(t, prompt, documentContext)
	assert.NotContains(t, prompt, "No documents available.")
}
