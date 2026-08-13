package llm

import (
	"fmt"
	"strings"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type ContextBuilder struct{}

func NewContextBuilder() *ContextBuilder {
	return &ContextBuilder{}
}

func (b *ContextBuilder) Build(
	question string,
	notes []models.Note,
	messages []models.ChatMessage,
	documentContext string,
) string {

	var builder strings.Builder

	builder.WriteString("SYSTEM INSTRUCTIONS:\n\n")

	builder.WriteString(
		"You are MindK AI, a knowledge management assistant. Answer the user's " +
			"question using only the KNOWLEDGE CONTEXT, DOCUMENT CONTEXT and " +
			"CONVERSATION HISTORY sections below - never external or general knowledge.\n\n",
	)

	builder.WriteString("Rules:\n")

	builder.WriteString("- Use only the notes, documents and conversation history shown below as your source of truth.\n")
	builder.WriteString("- Never assume a document, note or fact exists just because the user's question mentions or implies it - only trust what is explicitly shown below.\n")
	builder.WriteString("- If the answer is not explicitly available in the context, say plainly that the information is not available in the knowledge base. Do not guess, and do not fall back to general knowledge, even if you know the answer.\n")
	builder.WriteString("- When you do answer from the DOCUMENT CONTEXT, you may mention which document(s) it came from by name.\n")
	builder.WriteString("- Format responses using Markdown when useful.\n\n")

	builder.WriteString("CONVERSATION HISTORY:\n\n")

	if len(messages) == 0 {

		builder.WriteString("No previous conversation.\n\n")

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
				fmt.Sprintf("%s:\n%s\n\n", role, message.Content),
			)
		}
	}

	builder.WriteString("KNOWLEDGE CONTEXT:\n\n")

	if len(notes) == 0 {

		builder.WriteString("No notes available.\n\n")

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

	builder.WriteString("DOCUMENT CONTEXT:\n\n")

	if documentContext == "" {

		builder.WriteString("No documents available.\n\n")

	} else {

		builder.WriteString(documentContext)
		builder.WriteString("\n\n")
	}

	builder.WriteString("USER QUESTION:\n\n")
	builder.WriteString(question)
	builder.WriteString("\n\nRESPONSE:\n")

	return builder.String()
}
