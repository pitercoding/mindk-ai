package utils

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testFile satisfies multipart.File on top of a strings.Reader.
type testFile struct {
	*strings.Reader
}

func (testFile) Close() error { return nil }

// failingFile always fails to read, to exercise the error path.
type failingFile struct {
	testFile
}

func (failingFile) Read([]byte) (int, error) { return 0, errors.New("disk on fire") }

func TestReadFile(t *testing.T) {
	t.Run("reads markdown and text files as they are", func(t *testing.T) {
		for _, ext := range []string{".md", ".txt"} {
			got, err := ReadFile(testFile{strings.NewReader("# hello\nworld")}, ext)

			require.NoError(t, err, ext)
			assert.Equal(t, "# hello\nworld", got, ext)
		}
	})

	t.Run("rejects unsupported extensions", func(t *testing.T) {
		got, err := ReadFile(testFile{strings.NewReader("x")}, ".exe")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported file type: .exe")
		assert.Empty(t, got)
	})

	t.Run("wraps read errors", func(t *testing.T) {
		_, err := ReadFile(failingFile{testFile{strings.NewReader("x")}}, ".txt")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read uploaded file")
	})
}
