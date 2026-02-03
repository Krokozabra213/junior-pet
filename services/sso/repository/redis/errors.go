package redis

import "errors"

var (
	ErrInternal     = errors.New("internal error")
	ErrNotFound     = errors.New("not found error")
	ErrTokenExpired = errors.New("token already expired")
)
