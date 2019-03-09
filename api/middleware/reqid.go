package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/alextanhongpin/go-microservice/pkg/reqid"
)

// RequestID obtains the request id from the X-Request-Id header if present, or
// creates a new one and populates the context with it.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := reqid.FromHeader(c.Writer, c.Request)

		// Set in context.Context. This is preferable, since we can
		// just pass down the context to the next layer without
		// additional work.
		ctx := reqid.ContextWithRequestID(c.Request.Context(), reqID)
		c.Request = c.Request.WithContext(ctx)

		// Set in *gin.Context.
		// c.Set(requestIdKey, reqID)
		c.Next()
	}
}

// Alternative. Not preferred.
// const requestIdKey = "request_id"
// func GetRequestID(c *gin.Context) string {
//         v, _ := c.Get(requestIdKey)
//         reqID, _ := v.(string)
//         return reqID
// }
