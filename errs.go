// Package errs provides error codes and error handling.
package errs

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
)

// Common errors.
var (
	BadRequest          = New(CodeBadRequest, http.StatusText(http.StatusBadRequest))
	Unauthorized        = New(CodeUnauthorized, http.StatusText(http.StatusUnauthorized))
	Forbidden           = New(CodeForbidden, http.StatusText(http.StatusForbidden))
	NotFound            = New(CodeNotFound, http.StatusText(http.StatusNotFound))
	Gone                = New(CodeGone, http.StatusText(http.StatusGone))
	TooManyRequest      = New(CodeTooManyRequests, http.StatusText(http.StatusTooManyRequests))
	InternalServerError = New(CodeInternalServerError, http.StatusText(http.StatusInternalServerError))
	NotImplemented      = New(CodeNotImplemented, http.StatusText(http.StatusNotImplemented))
	ServiceUnavailable  = New(CodeServiceUnavailable, http.StatusText(http.StatusServiceUnavailable))
)

// Code represents an error code.
type Code string

// String returns the string representation of the error code.
func (c Code) String() string {
	return string(c)
}

// Error codes.
const (
	CodeBadRequest      Code = "BAD_REQUEST"
	CodeUnauthorized    Code = "UNAUTHORIZED"
	CodeForbidden       Code = "FORBIDDEN"
	CodeNotFound        Code = "NOT_FOUND"
	CodeGone            Code = "GONE"
	CodeTooManyRequests Code = "TOO_MANY_REQUESTS"

	CodeInternalServerError Code = "INTERNAL_SERVER_ERROR"
	CodeNotImplemented      Code = "NOT_IMPLEMENTED"
	CodeServiceUnavailable  Code = "SERVICE_UNAVAILABLE"
)

// Error represents an error.
type Error struct {
	// Code is the error code.
	Code Code `json:"code"`

	// Message is the error message.
	Message string `json:"message"`

	// Info is additional information about the error.
	Info map[string]interface{} `json:"info,omitempty"`

	// Timestamp is the time when the error occurred.
	Timestamp time.Time `json:"timestamp"`
}

// Error returns the string representation of the error.
func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// HTTPStatusCode returns the HTTP status code for the error.
func (e *Error) HTTPStatusCode() int {
	switch e.Code {
	case CodeBadRequest:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeGone:
		return http.StatusGone
	case CodeTooManyRequests:
		return http.StatusTooManyRequests
	case CodeInternalServerError:
		return http.StatusInternalServerError
	case CodeNotImplemented:
		return http.StatusNotImplemented
	case CodeServiceUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// Option represents an option for an error.
type Option func(*option)

// option represents an option.
type option struct {
	info   map[string]interface{}
	logErr error
}

// WithInfo sets the info option.
func WithInfo(info map[string]interface{}) Option {
	return func(o *option) {
		o.info = info
	}
}

// WithLogErr sets the log error option.
func WithLogErr(err error) Option {
	return func(o *option) {
		o.logErr = err
	}
}

// New returns a new error.
func New(code Code, msg string, opts ...Option) *Error {
	o := new(option)
	for _, opt := range opts {
		opt(o)
	}

	if o.logErr != nil {
		logrus.WithError(o.logErr).Error(msg)
	}

	e := &Error{
		Code:      code,
		Message:   msg,
		Timestamp: time.Now(),
		Info:      o.info,
	}

	return e
}

// InvalidStructError returns a new error for an invalid struct.
func InvalidStructError(err error) *Error {
	return New(CodeBadRequest, http.StatusText(http.StatusBadRequest), WithInfo(validationInfo(err)))
}

// validationInfo returns the validation info for the error.
func validationInfo(err error) map[string]interface{} {
	result := make(map[string]interface{})
	if errCast, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errCast {
			result[strcase.ToLowerCamel(e.Field())] = toMessage(e)
		}

		return result
	}

	result["error"] = err.Error()
	return result
}

// toMessage returns the message for the validation error.
func toMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", strcase.ToLowerCamel(e.Field()))
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", strcase.ToLowerCamel(e.Field()), e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s", strcase.ToLowerCamel(e.Field()), e.Param())
	case "email":
		return "invalid email format"
	case "len":
		return fmt.Sprintf("%s must be %s characters long", strcase.ToLowerCamel(e.Field()), e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be %s", strcase.ToLowerCamel(e.Field()), e.Param())
	}

	return fmt.Sprintf("%s is not valid", strcase.ToLowerCamel(e.Field()))
}

// ResponseError returns an error response.
func ResponseError(c *gin.Context, err error) {
	var e *Error
	if ok := errors.As(err, &e); ok {
		c.JSON(e.HTTPStatusCode(), e)
		return
	}

	c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}
