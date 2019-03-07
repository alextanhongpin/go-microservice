package model

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/alextanhongpin/go-microservice/pkg/reqid"
)

// ErrorResponse represents a json error response.
type ErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"error,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

func NewErrorResponse(c *gin.Context, err error) *ErrorResponse {
	// Get the request id and return it in the error response. This allows
	// us to trace the error by allowing the user (client-facing) to submit
	// the returned code to ops when reporting the error.
	ctx := c.Request.Context()
	reqID, _ := reqid.FromContext(ctx)
	return &ErrorResponse{
		Code:    reqID,
		Message: err.Error(),
	}
}

// ErrorJSON returns a basic error json with the error code.
func ErrorJSON(c *gin.Context, err error) {

	// TODO: Set the error in the gin context too - this allows us to centralize
	// the error logging.
	c.Error(err)
	c.JSON(http.StatusBadRequest, NewErrorResponse(c, err))
}
