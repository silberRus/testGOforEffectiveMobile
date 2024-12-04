package errors

import "fmt"

type ErrorType string

const (
	NotFound      ErrorType = "NOT_FOUND"
	BadRequest    ErrorType = "BAD_REQUEST"
	Internal      ErrorType = "INTERNAL"
	Validation    ErrorType = "VALIDATION"
	AlreadyExists ErrorType = "ALREADY_EXISTS"
)

type Error struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewNotFound(message string, err error) *Error {
	return &Error{
		Type:    NotFound,
		Message: message,
		Err:     err,
	}
}

func NewLyricsNotFound(message string, err error) *Error {
	return &Error{
		Type:    NotFound,
		Message: message,
		Err:     err,
	}
}

func NewBadRequest(message string, err error) *Error {
	return &Error{
		Type:    BadRequest,
		Message: message,
		Err:     err,
	}
}

func NewInternal(message string, err error) *Error {
	return &Error{
		Type:    Internal,
		Message: message,
		Err:     err,
	}
}

func NewValidation(message string, err error) *Error {
	return &Error{
		Type:    Validation,
		Message: message,
		Err:     err,
	}
}

func NewAlreadyExists(message string, err error) *Error {
	return &Error{
		Type:    AlreadyExists,
		Message: message,
		Err:     err,
	}
}
