package xreqid

import (
	"context"
	"net/http"

	"github.com/rs/xid"
)

type contextKey string

const requestIdContextKey = "req_id"

const HeaderXRequestID = "X-Request-Id"

func New() string {
	return xid.New().String()
}

func NewContext() context.Context {
	ctx := context.Background()
	reqID := New()
	return ContextWithRequestID(ctx, reqID)
}

func FromHeader(w http.ResponseWriter, r *http.Request) string {
	reqID := r.Header.Get(HeaderXRequestID)
	if reqID == "" {
		reqID = New()
		// Set the header for both the current request and also
		// response.
		r.Header.Set(HeaderXRequestID, reqID)
		w.Header().Set(HeaderXRequestID, reqID)
	}
	return reqID
}

func ContextWithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, requestIdContextKey, reqID)
}

func FromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIdContextKey).(string)
	return id, ok
}
