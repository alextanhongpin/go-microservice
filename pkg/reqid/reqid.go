package reqid

import (
	"context"
	"net/http"

	"github.com/rs/xid"
)

type contextKey string

const requestIDContextKey = contextKey("req_id")

// NOTE: The ID is uppercase.
const headerXRequestID = "X-Request-ID"

// New returns a unique request id.
func New() string {
	return xid.New().String()
}

// FromHeader attempts to obtain a request id from the X-Request-Id header, or
// creates a new one if it does not exist.
func FromHeader(w http.ResponseWriter, r *http.Request) string {
	reqID := r.Header.Get(headerXRequestID)
	if reqID == "" {
		reqID = New()
		// Set the header for both the current request and also
		// response.
		r.Header.Set(headerXRequestID, reqID)
		w.Header().Set(headerXRequestID, reqID)
	}
	return reqID
}

func ContextWithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey, reqID)
}

func FromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDContextKey).(string)
	return id, ok
}
