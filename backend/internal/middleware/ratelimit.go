package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/pitercoding/mindk-ai/backend/internal/auth"
)

// visitor tracks one authenticated user's request count within the current
// fixed window.
type visitor struct {
	count       int
	windowStart time.Time
}

// RateLimiter enforces a fixed-window request limit per authenticated user.
// It only knows how many requests a user has made and when their window
// started - it has no knowledge of notes, documents or chat.
type RateLimiter struct {
	limit  int
	window time.Duration
	now    func() time.Time

	mu       sync.Mutex
	visitors map[string]*visitor
}

// NewRateLimiter creates a rate limiter allowing up to limit requests per
// window, tracked independently per authenticated user.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:    limit,
		window:   window,
		now:      time.Now,
		visitors: make(map[string]*visitor),
	}
}

// Middleware rejects requests exceeding the limit with
// 429 Too Many Requests and a Retry-After header. It identifies callers by
// the authenticated user ID set in the request context, so it must run
// after the auth middleware.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		allowed, retryAfter := rl.allow(userID)
		if !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// allow records a request for userID and reports whether it is within the
// current window's limit. When it isn't, it also returns how long until the
// window resets.
func (rl *RateLimiter) allow(userID string) (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := rl.now()

	v, exists := rl.visitors[userID]
	if !exists || now.Sub(v.windowStart) >= rl.window {
		rl.visitors[userID] = &visitor{count: 1, windowStart: now}
		return true, 0
	}

	if v.count >= rl.limit {
		return false, rl.window - now.Sub(v.windowStart)
	}

	v.count++
	return true, 0
}
