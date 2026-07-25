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

	builder.WriteString(
		"SYSTEM INSTRUCTIONS:\n\n",
	)

	builder.WriteString(
		"You are MindK AI, an assistant that answers questions using provided knowledge.\n\n",
	)

	builder.WriteString(
		"Rules:\n",
	)

	builder.WriteString(
		"- Use only the provided notes and conversation history to answer.\n",
	)

	builder.WriteString(
		"- Do not use external knowledge or assumptions.\n",
	)

	builder.WriteString(
		"-  If the answer is not explicitly available in the context, say that there is not enough information.\n",
	)

	builder.WriteString(
		"- Format responses using Markdown when useful.\n\n",
	)

	builder.WriteString(
		"CONVERSATION HISTORY:\n\n",
	)

	if len(messages) == 0 {

		builder.WriteString(
			"No previous conversation.\n\n",
		)

	} else {

		for _, message := range messages {

			role := message.Role

			switch role {
			case "user":
				role = "User"

			case "assistant":
				role = "Assistant"
			}

			builder.WriteString(
				fmt.Sprintf(
					"%s:\n%s\n\n",
					role,
					message.Content,
				),
			)
		}
	}

	builder.WriteString(
		"KNOWLEDGE CONTEXT:\n\n",
	)

	if len(notes) == 0 {

		builder.WriteString(
			"No notes available.\n\n",
		)

	} else {

		for _, note := range notes {

			builder.WriteString(
				fmt.Sprintf(
					"Title: %s\nContent:\n%s\n\n",
					note.Title,
					note.Content,
				),
			)
		}
	}

	builder.WriteString(
		"USER QUESTION:\n\n",
	)

	builder.WriteString(question)

	builder.WriteString(
		"\n\nRESPONSE:\n",
	)

	return builder.String()
}
