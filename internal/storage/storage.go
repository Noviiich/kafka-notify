package storage

import (
	"errors"
)

var ErrUserNotFound = errors.New("user not found")
var ErrNoNotificationsFound = errors.New("no notifications found for user")
