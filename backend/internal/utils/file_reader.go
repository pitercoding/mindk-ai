package utils

import (
	"fmt"
	"io"
	"mime/multipart"
)

func ReadFile(
	file multipart.File,
	extension string,
) (string, error) {

	switch extension {

		case ".pdf":
			return ReadPDF(file)

		case ".md", ".txt":
			bytes, err := io.ReadAll(file)

			if err != nil {
				return "", fmt.Errorf(
					"failed to read uploaded file: %w",
					err,
				)
			}

			return string(bytes), nil

		default:
			return "", fmt.Errorf(
				"unsupported file type: %s",
				extension,
			)
	}
}
