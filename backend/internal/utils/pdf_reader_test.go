package utils

import (
	"os"
	"testing"
)

func TestReadPDF(t *testing.T) {

	file, err := os.Open("cv.pdf")

	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()

	content, err := ReadPDF(file)

	if err != nil {
		t.Fatal(err)
	}

	if content == "" {
		t.Fatal("pdf content is empty")
	}

	t.Log(content[:100])
}
