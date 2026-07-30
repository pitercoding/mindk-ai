package utils

import "strings"

// SplitIntoChunks divides a text into chunks based on paragraphs.
func SplitIntoChunks(text string, chunkSize int) []string {

	text = strings.TrimSpace(text)

	if text == "" {
		return []string{}
	}

	if chunkSize <= 0 {
		chunkSize = 1
	}

	paragraphs := strings.Split(text, "\n\n")

	chunks := make([]string, 0)

	for i := 0; i < len(paragraphs); i += chunkSize {

		end := i + chunkSize

		if end > len(paragraphs) {
			end = len(paragraphs)
		}

		chunk := strings.Join(paragraphs[i:end], "\n\n")

		chunks = append(chunks, chunk)
	}

	return chunks
}
