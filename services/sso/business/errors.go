package business

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInternal         = errors.New("internal service error")
	ErrPermissionDenied = errors.New("permission denied")
	ErrUserExists       = errors.New("user already exists")
	ErrEmailExists      = errors.New("user email already exists")
)
