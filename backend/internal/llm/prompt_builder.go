package llm

import (
	"fmt"
	"strings"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

func BuildPrompt(
	question string,
	notes []models.Note,
	messages []models.ChatMessage,
) string {

	var builder strings.Builder

	builder.WriteString("You are an AI assistant.\n\n")

	builder.WriteString("Use the conversation history and notes to answer.\n\n")

	builder.WriteString("CONVERSATION HISTORY:\n\n")

	for _, message := range messages {

		role := message.Role

		if role == "user" {
			role = "User"
		} else if role == "assistant" {
			role = "Assistant"
		}

		builder.WriteString(fmt.Sprintf(
			"%s: %s\n",
			role,
			message.Content,
		))
	}

	builder.WriteString("\n")

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
