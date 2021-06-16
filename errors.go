package bubucore

import (
	"net/http"
)

// Errors
var (
	ErrUserBlocked  = NewError(http.StatusForbidden, "requested user is blocked")
	ErrTokenInvalid = NewError(http.StatusUnauthorized, "token is invalid")
	ErrTokenExpired = NewError(http.StatusUnauthorized, "token is expired")
)

// NewError creates a new Error instance
func NewError(code int, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

// Error defines the response error
type Error struct {
	Code    int    `json:"code" example:"403"`
	Message string `json:"message" example:"Access denied"`
}

// Error as a string
func (e *Error) Error() string {
	return e.Message
}
