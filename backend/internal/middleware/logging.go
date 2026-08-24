package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"github.com/pitercoding/mindk-ai/backend/internal/auth"
)

type logFieldsKey struct{}

// logFields is shared, via a pointer stashed in the request context, between
// RequestLogger and any middleware/handler that runs downstream of it (e.g.
// Clerk auth). Downstream code fills in fields as they become known; when
// the handler returns, RequestLogger reads the final values back out.
type logFields struct {
	userID string
}

// newRequestID returns a short random hex correlation ID, unique enough to
// tell requests apart in logs without needing a full UUID dependency.
func newRequestID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// recordUserID attaches userID to the current request's access log, if
// RequestLogger is running for this request. It is a no-op otherwise (e.g.
// in handler unit tests that call handlers directly).
func recordUserID(ctx context.Context, userID string) {
	if lf, ok := ctx.Value(logFieldsKey{}).(*logFields); ok {
		lf.userID = userID
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// RequestLogger assigns each request a correlation ID (returned to the
// client via the X-Request-ID header and available downstream via
// auth.RequestIDFromContext), then logs one structured line per request with
// method, path, status, duration and - when the request is authenticated -
// the user ID. It never logs request/response bodies.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := newRequestID()
		w.Header().Set("X-Request-ID", requestID)

		lf := &logFields{}
		ctx := auth.WithRequestID(r.Context(), requestID)
		ctx = context.WithValue(ctx, logFieldsKey{}, lf)
		r = r.WithContext(ctx)

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		attrs := []any{
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
			"request_id", requestID,
		}
		if lf.userID != "" {
			attrs = append(attrs, "user_id", lf.userID)
		}

		slog.Info("request completed", attrs...)
	})
}
