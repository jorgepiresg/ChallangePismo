package utils

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Error struct {
	HTTPCode int         `mapstructure:"code" json:"-"`
	Message  string      `mapstructure:"message" json:"message,omitempty"`
	Detail   interface{} `mapstructure:"detail,omitempty" swaggerignore:"true" json:"detail,omitempty"`
}

func NewError(httpCode int, message string, detail interface{}) *Error {
	return &Error{
		HTTPCode: httpCode,
		Message:  message,
		Detail:   detail,
	}
}

func NewErrorContext(c echo.Context, httpCode int, message string, detail interface{}) *Error {
	c.NoContent(httpCode)
	return &Error{
		HTTPCode: httpCode,
		Message:  message,
		Detail:   detail,
	}
}

func GetError(err error) *Error {
	if err == nil {
		return nil
	}

	e, ok := err.(*Error)
	if !ok {
		return &Error{
			HTTPCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}
	}

	return e
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %v - message: %v - detail: %v", e.HTTPCode, e.Message, string(ToJSON(e.Detail)))
}
