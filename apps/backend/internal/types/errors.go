package types

import "github.com/gin-gonic/gin"

type APIError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Cause   string `json:"cause"`
	Message string `json:"message"`
}

func (e *APIError) Respond(c *gin.Context) {
	c.AbortWithStatusJSON(e.Status, e)
}

func NewAPIError(status int, code, message string, cause error) *APIError {
	var causeMess string
	if cause != nil {
		causeMess = cause.Error()
	}
	return &APIError{
		Status:  status,
		Code:    code,
		Cause:   causeMess,
		Message: message,
	}
}
