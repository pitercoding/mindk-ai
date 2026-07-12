package llm

import (
	"fmt"
	"strings"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

func BuildPrompt(question string, notes []models.Note, history []models.ChatHistory) string {

	_ = history
	
	var builder strings.Builder

	builder.WriteString("You are an AI assistant.\n\n")

	builder.WriteString("Answer ONLY using the notes below.\n")

	builder.WriteString("If the answer is not present, say that the information is not available.\n\n")

	builder.WriteString("NOTES:\n\n")

	for _, note := range notes {
		builder.WriteString(fmt.Sprintf(
			"Title: %s\nContent: %s\n\n",
			note.Title,
			note.Content,
		))
	}

	builder.WriteString(fmt.Sprintf(
		"USER QUESTION:\n%s",
		question,
	))

	return builder.String()
}
