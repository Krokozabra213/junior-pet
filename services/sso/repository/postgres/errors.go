package postgres

import "errors"

var (
	ErrInternal          = errors.New("internal error")
	ErrNotFound          = errors.New("not found error")
	ErrUserAlreadyExists = errors.New("user already exists")
)
