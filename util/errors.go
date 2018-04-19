package util

import "errors"

var (
	// ErrNotFound to indicate that a required subject was not found
	ErrNotFound = errors.New("not found")
)

// baseError is a trivial implementation of error.
type baseError struct {
	s string
}

func (e *baseError) Error() string {
	return e.s
}

type validationError struct {
	baseError
}
