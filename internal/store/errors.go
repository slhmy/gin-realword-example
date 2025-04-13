package store

import (
	"errors"
	"net/http"

	gin_utils "gin-realword-example/internal/modules/utils/gin"
)

var (
	ErrStoreCreateUserFailed     = errors.New("store: user create failed")
	ErrStoreUserNotFoundByID     = errors.New("store: failed to get user by id")
	ErrStoreUserNotFoundByEmail  = errors.New("store: failed to get user by email")
	ErrStoreLoginSessionNotFound = errors.New("store: failed to get login session")
)

func init() {
	gin_utils.RegisterErrHttpStatusMapping(ErrStoreUserNotFoundByID, http.StatusNotFound)
	gin_utils.RegisterErrHttpStatusMapping(ErrStoreUserNotFoundByEmail, http.StatusNotFound)
	gin_utils.RegisterErrHttpStatusMapping(ErrStoreLoginSessionNotFound, http.StatusUnauthorized)
}
