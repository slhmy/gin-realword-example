package services

import (
	"context"
	"errors"
	"gin-realword-example/internal/models"
	"gin-realword-example/internal/store"
	"net/http"

	gwm_app "github.com/slhmy/go-webmods/app"
)

func GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := store.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrStoreUserNotFoundByID) {
			return nil, gwm_app.ServiceError{Err: err, HttpStatusCode: http.StatusNotFound}
		}
	}
	return user, nil
}
