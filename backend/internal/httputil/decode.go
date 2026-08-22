package httputil

import (
	"encoding/json"
	"errors"
	"net/http"
)

// MaxJSONBodyBytes is the default body size limit for JSON endpoints
// (notes, chat, chat messages, chat sessions). It is generous for
// conversational/text content while still ruling out multi-megabyte
// bodies aimed at exhausting memory.
const MaxJSONBodyBytes int64 = 1 << 20 // 1 MiB

// DecodeJSON reads at most maxBytes from r.Body and decodes it as JSON into
// v. On failure it writes the appropriate error response itself - 413 if the
// body exceeded maxBytes, 400 for any other decode error - and returns a
// non-nil error so the caller can just return.
func DecodeJSON(w http.ResponseWriter, r *http.Request, maxBytes int64, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return err
		}

		http.Error(w, "invalid request body", http.StatusBadRequest)
		return err
	}

	return nil
}
