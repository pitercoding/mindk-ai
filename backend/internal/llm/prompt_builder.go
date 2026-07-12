package llm

import (
	"fmt"
	"strings"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

func BuildPrompt(
	question string,
	notes []models.Note,
	history []models.ChatHistory,
) string {

	var builder strings.Builder

	builder.WriteString("You are an AI assistant.\n\n")

	builder.WriteString("Use the conversation history and notes to answer.\n\n")

	builder.WriteString("CONVERSATION HISTORY:\n\n")

	for _, item := range history {

		builder.WriteString(fmt.Sprintf(
			"User: %s\nAssistant: %s\n\n",
			item.Question,
			item.Answer,
		))
	}

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
