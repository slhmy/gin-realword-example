package services

import (
	"context"
	"errors"
	"gin-realword-example/internal/models"
	"gin-realword-example/internal/store"

	gwm_app "github.com/slhmy/go-webmods/app"
)

func GetUserByID(ctx context.Context, id gwm_app.ID) (*models.User, error) {
	user, err := store.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrStoreUserNotFoundByID) {
			return nil, gwm_app.NewNotFoundError("user", "id")
		}
	}
	return user, nil
}
