package httputil

import (
	"net/http"
	"strconv"
	"strings"
)

func GetIDFromPath(r *http.Request) (int, error) {
	parts := strings.Split(r.URL.Path, "/")

	idStr := parts[len(parts)-1]

	return strconv.Atoi(idStr)
}
