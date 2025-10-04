package domain

import "errors"

var ErrInvalid = errors.New("invalid")
var ErrNotFount = errors.New("not found")
var ErrConflict = errors.New("conflict")
var ErrToLarge = errors.New("too large")

type InvalidFieldError struct {
	Field string
	Reason string
}

func (e *InvalidFieldError) Error() string {
	return e.Reason
}