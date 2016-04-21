package errutil

import (
	"errors"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidParameter = errors.New("invalid parametr")
	ErrAlreadyExisted   = errors.New("already existed")
)
