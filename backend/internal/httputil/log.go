package httputil

import (
	"log/slog"
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/auth"
)

// logAttrs builds the request_id/user_id attributes shared by every log line
// tied to a request, pulling them from context rather than requiring callers
// to thread them through by hand.
func logAttrs(r *http.Request) []any {
	var attrs []any

	if requestID, ok := auth.RequestIDFromContext(r.Context()); ok {
		attrs = append(attrs, "request_id", requestID)
	}
	if userID, ok := auth.UserIDFromContext(r.Context()); ok {
		attrs = append(attrs, "user_id", userID)
	}

	return attrs
}

// LogError logs an internal error tied to the current request. Callers must
// only pass the underlying error - never request bodies, note/document/chat
// content, or credentials - since it is written as-is.
func LogError(r *http.Request, message string, err error) {
	attrs := append([]any{"error", err}, logAttrs(r)...)
	slog.Error(message, attrs...)
}

// LogInfo logs a notable, non-error event tied to the current request (e.g.
// a completed upload). attrs must never include request bodies or content.
func LogInfo(r *http.Request, message string, attrs ...any) {
	slog.Info(message, append(attrs, logAttrs(r)...)...)
}
