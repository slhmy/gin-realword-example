package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gin-realword-example/internal/models"
	redis_client "gin-realword-example/internal/modules/clients/redis"
	"gin-realword-example/internal/modules/core"
	"gin-realword-example/internal/modules/shared"
	"gin-realword-example/internal/modules/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"

	gwm_app "github.com/slhmy/go-webmods/app"
	gwm_gorm "github.com/slhmy/go-webmods/clients/gorm"
)

var loginSessionExpireIn time.Duration

func init() {
	loginSessionExpireIn = core.ConfigStore.GetDuration(shared.ConfigKeyAuthLoginSessionExpireIn)
	err := gwm_gorm.GetDB().AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}
}

func UpsertUser(ctx context.Context, request models.CreateUserRequest) (*gwm_app.ID, error) {
	_, err := GetUserByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, ErrStoreUserNotFoundByEmail) {
			return CreateUser(ctx, request)
		}
		return nil, err
	}
	return nil, nil
}

func CreateUser(ctx context.Context, request models.CreateUserRequest) (*gwm_app.ID, error) {
	user := models.User{
		Username: request.Username,
		Email:    request.Email,
	}
	err := gwm_gorm.GetDB().Create(&user).Error
	if err != nil {
		return nil, ErrStoreCreateUserFailed
	}
	return &user.ID, nil
}

func CheckUserExistsByID(ctx context.Context, id gwm_app.ID) (bool, error) {
	var count int64
	err := gwm_gorm.GetDB().Table("users").Where("id = ?", id).Count(&count).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return count > 0, nil
}

func GetUserByID(ctx context.Context, id gwm_app.ID) (*models.User, error) {
	var user models.User
	err := gwm_gorm.GetDB().First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Join(ErrStoreUserNotFoundByID, err)
		}
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := gwm_gorm.GetDB().First(&user, "email = ?", email).Error
	if err != nil {
		return nil, errors.Join(ErrStoreUserNotFoundByEmail, err)
	}
	return &user, nil
}

func GenerateLoginSession(ctx context.Context, id gwm_app.ID) (*string, *time.Time, error) {
	exist, err := CheckUserExistsByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if !exist {
		return nil, nil, ErrStoreUserNotFoundByID
	}
	sessionID := uuid.New().String()
	expireAt := time.Now().Add(loginSessionExpireIn)
	err = redis_client.GetRDB().
		SetEx(ctx, fmt.Sprintf(AuthKeyFormatLoginSession, sessionID), string(id), loginSessionExpireIn).
		Err()
	if err != nil {
		return nil, nil, err
	}
	return &sessionID, &expireAt, nil
}

func GetUserIDFromLoginSession(ctx context.Context, sessionID string) (*string, *time.Time, error) {
	sessionKey := fmt.Sprintf(AuthKeyFormatLoginSession, sessionID)
	result, err := redis_client.GetRDB().Get(ctx, sessionKey).Result()
	if err != nil {
		return nil, nil, ErrStoreLoginSessionNotFound
	}
	expireAt := time.Now().Add(loginSessionExpireIn)
	err = redis_client.GetRDB().
		SetEx(ctx, sessionKey, result, loginSessionExpireIn).
		Err()
	if err != nil {
		return nil, nil, err
	}
	return utils.ToPtr(result), &expireAt, nil
}
