package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type envelope struct {
	Data  any            `json:"data"`
	Error *apiError      `json:"error"`
	Meta  map[string]any `json:"meta"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func ok(c *gin.Context, data any, meta map[string]any) {
	if meta == nil {
		meta = map[string]any{}
	}
	c.JSON(http.StatusOK, envelope{Data: data, Error: nil, Meta: meta})
}

func created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, envelope{Data: data, Error: nil, Meta: map[string]any{}})
}

func fail(c *gin.Context, status int, code, message string, details any) {
	c.JSON(status, envelope{Data: nil, Error: &apiError{Code: code, Message: message, Details: details}, Meta: map[string]any{}})
}
