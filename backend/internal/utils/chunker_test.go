package utils

import (
	"strings"
	"testing"
)

func TestSplitIntoChunks(t *testing.T) {

	t.Run("should return empty slice for empty text", func(t *testing.T) {

		chunks := SplitIntoChunks("", 500)

		if len(chunks) != 0 {
			t.Errorf(
				"expected 0 chunks, got %d",
				len(chunks),
			)
		}
	})

	t.Run("should return one chunk when text is smaller than chunk size", func(t *testing.T) {

		text := "Go is a simple programming language."

		chunks := SplitIntoChunks(text, 500)

		if len(chunks) != 1 {
			t.Fatalf(
				"expected 1 chunk, got %d",
				len(chunks),
			)
		}

		if chunks[0] != text {
			t.Errorf(
				"expected %q, got %q",
				text,
				chunks[0],
			)
		}
	})

	t.Run("should split long text into multiple chunks", func(t *testing.T) {

		text := strings.Repeat(
			"golang ",
			200,
		)

		chunks := SplitIntoChunks(text, 500)

		if len(chunks) < 2 {
			t.Fatalf(
				"expected multiple chunks, got %d",
				len(chunks),
			)
		}
	})

	t.Run("should not create chunks larger than chunk size", func(t *testing.T) {

		text := strings.Repeat(
			"golang ",
			200,
		)

		chunks := SplitIntoChunks(text, 500)

		for i, chunk := range chunks {

			if len(chunk) > 500 {
				t.Errorf(
					"chunk %d has %d characters",
					i,
					len(chunk),
				)
			}
		}
	})

	t.Run("should not lose content after splitting", func(t *testing.T) {

		text := strings.Repeat(
			"golang ",
			200,
		)

		chunks := SplitIntoChunks(text, 500)

		result := strings.Join(
			chunks,
			" ",
		)

		expected := strings.Join(
			strings.Fields(text),
			" ",
		)

		if result != expected {
			t.Error(
				"text changed after chunking",
			)
		}
	})
}
