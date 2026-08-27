package middleware

import (
	"net/http"
	"strings"
)

// swaggerCSP relaxes just enough of the default-src:none policy for the
// Swagger UI page (inline <script>/<style> and same-origin JS/CSS assets
// served by http-swagger) to render. It is only ever reached in
// development, since /swagger/ is not registered in production.
const swaggerCSP = "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self' data:"

// SecurityHeaders sets baseline security headers on every response. It is
// pure API middleware: MindK AI's backend never renders HTML, so the policy
// is deliberately locked down rather than tuned for embedding or scripts -
// except for the development-only Swagger UI, which needs a less strict CSP
// to render at all.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Cache-Control", "no-store")

		if strings.HasPrefix(r.URL.Path, "/swagger/") {
			h.Set("Content-Security-Policy", swaggerCSP)
		} else {
			h.Set("Content-Security-Policy", "default-src 'none'")
		}

		next.ServeHTTP(w, r)
	})
}
