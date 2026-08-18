package middleware

import (
	"net/http"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
)

// clockLeeway absorbs small clock drift between this server and Clerk's,
// which would otherwise fail verification of tokens that are technically
// valid (e.g. "token issued in the future").
const clockLeeway = 5 * time.Second

// NewClerkAuth configures the Clerk SDK with the given secret key and
// returns a middleware that verifies the Clerk session token sent in
// the Authorization header of incoming requests.
//
// Requests without a valid session are rejected with 401 Unauthorized.
// On success, the Clerk SessionClaims are available in the request
// context via clerk.SessionClaimsFromContext.
func NewClerkAuth(secretKey string) func(http.Handler) http.Handler {
	clerk.SetKey(secretKey)

	return func(next http.Handler) http.Handler {
		return clerkhttp.WithHeaderAuthorization(clerkhttp.Leeway(clockLeeway))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if !ok || claims == nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}))
	}
}
