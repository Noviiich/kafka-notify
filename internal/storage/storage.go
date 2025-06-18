package storage

import (
	"errors"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotUpdated = errors.New("user not updated")
