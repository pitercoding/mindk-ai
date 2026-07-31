package utils

import "strings"

// SplitIntoChunks divides text into chunks based on character count
// without breaking words.
func SplitIntoChunks(
	text string,
	chunkSize int,
) []string {

	text = strings.TrimSpace(text)

	if text == "" {
		return []string{}
	}

	if chunkSize <= 0 {
		chunkSize = 500
	}

	words := strings.Fields(text)

	chunks := make([]string, 0)

	var builder strings.Builder

	currentSize := 0

	for _, word := range words {

		wordSize := len(word)

		additional := wordSize

		if currentSize > 0 {
			additional++
		}

		if currentSize+additional > chunkSize {

			chunks = append(
				chunks,
				builder.String(),
			)

			builder.Reset()
			currentSize = 0
		}

		if currentSize > 0 {
			builder.WriteByte(' ')
			currentSize++
		}

		builder.WriteString(word)
		currentSize += wordSize
	}

	if builder.Len() > 0 {
		chunks = append(
			chunks,
			builder.String(),
		)
	}

	return chunks
}
