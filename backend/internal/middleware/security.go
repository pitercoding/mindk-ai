package middleware

import "net/http"

// SecurityHeaders sets baseline security headers on every response. It is
// pure API middleware: MindK AI's backend never renders HTML, so the policy
// is deliberately locked down rather than tuned for embedding or scripts.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Content-Security-Policy", "default-src 'none'")
		h.Set("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}
