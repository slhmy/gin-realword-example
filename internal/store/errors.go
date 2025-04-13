package store

import (
	"errors"
)

var (
	ErrStoreCreateUserFailed     = errors.New("store: user create failed")
	ErrStoreUserNotFoundByID     = errors.New("store: failed to get user by id")
	ErrStoreUserNotFoundByEmail  = errors.New("store: failed to get user by email")
	ErrStoreLoginSessionNotFound = errors.New("store: failed to get login session")
)
