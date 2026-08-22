package httputil

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func GetIDFromPath(r *http.Request) (int, error) {
	parts := strings.Split(r.URL.Path, "/")

	idStr := parts[len(parts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	if id <= 0 {
		return 0, fmt.Errorf("id must be a positive integer, got %d", id)
	}

	return id, nil
}
