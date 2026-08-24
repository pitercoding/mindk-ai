package auth

import "context"

type contextKey int

const (
	userIDKey contextKey = iota
	requestIDKey
)

// WithUserID returns a copy of ctx carrying the authenticated user's ID.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext returns the authenticated user's ID stored in ctx by
// the auth middleware, and whether one was present.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

// WithRequestID returns a copy of ctx carrying the current request's
// correlation ID.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// RequestIDFromContext returns the correlation ID stored in ctx by the
// request logging middleware, and whether one was present.
func RequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDKey).(string)
	return requestID, ok
}
