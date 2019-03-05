package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/alextanhongpin/logging/pkg/xreqid"
)

// const requestIdKey = "request_id"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := xreqid.FromHeader(c.Writer, c.Request)

		// Set in context.Context. This is preferable, since we can
		// just pass down the context to the next layer without
		// additional work.
		ctx := xreqid.ContextWithRequestID(c.Request.Context(), reqID)
		c.Request = c.Request.WithContext(ctx)

		// Set in *gin.Context.
		// c.Set(requestIdKey, reqID)
		c.Next()
	}
}

// Alternative. Not preferred.
// func GetRequestID(c *gin.Context) string {
//         v, _ := c.Get(requestIdKey)
//         reqID, _ := v.(string)
//         return reqID
// }
