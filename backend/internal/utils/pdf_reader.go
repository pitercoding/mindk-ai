package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/ledongthuc/pdf"
)

func ReadPDF(file multipart.File) (string, error) {

	tempFile, err := os.CreateTemp("", "*.pdf")
	if err != nil {
		return "", fmt.Errorf(
			"failed to create temporary file: %w",
			err,
		)
	}

	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, file)
	if err != nil {
		tempFile.Close()

		return "", fmt.Errorf(
			"failed to save temporary pdf: %w",
			err,
		)
	}

	err = tempFile.Close()
	if err != nil {
		return "", fmt.Errorf(
			"failed to close temporary file: %w",
			err,
		)
	}

	pdfFile, reader, err := pdf.Open(tempFile.Name())

	if err != nil {
		return "", fmt.Errorf(
			"failed to open pdf: %w",
			err,
		)
	}

	defer pdfFile.Close()

	var content string

	totalPages := reader.NumPage()

	for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {

		page := reader.Page(pageIndex)

		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)

		if err != nil {
			return "", fmt.Errorf(
				"failed to extract pdf text: %w",
				err,
			)
		}

		content += text + "\n"
	}

	return content, nil
}
