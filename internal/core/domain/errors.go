package domain

import "errors"

var (
	ErrInternalError = errors.New("internal server error")
	ErrUserNotFound  = errors.New("user does not exist")
)
