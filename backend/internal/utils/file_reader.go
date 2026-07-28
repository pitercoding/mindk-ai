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

	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf(
			"failed to read uploaded file: %w",
			err,
		)
	}

	return string(bytes), nil
}
