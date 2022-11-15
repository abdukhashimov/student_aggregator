package domain

import "errors"

var (
	ErrInternalError = errors.New("internal server error")
	ErrNotFound      = errors.New("resource does not exist")
)
