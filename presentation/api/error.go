package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/alextanhongpin/pkg/requestid"
)

// Error represents a json error response.
type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"error,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

// NewError returns a new JSON error with the request id for debugging.
func NewError(c *gin.Context, err error) *Error {
	// Get the request id and return it in the error response. This allows
	// us to trace the error by allowing the user (client-facing) to submit
	// the returned code to ops when reporting the error.
	ctx := c.Request.Context()
	reqID, _ := requestid.Value(ctx)
	return &Error{
		Code:    reqID,
		Message: err.Error(),
	}
}

// ErrorJSON returns a basic error json with the error code.
func ErrorJSON(c *gin.Context, err error) {
	c.Error(err)
	c.JSON(http.StatusBadRequest, NewError(c, err))
}
