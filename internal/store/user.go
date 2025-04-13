package store

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gin-realword-example/internal/models"
	gorm_client "gin-realword-example/internal/modules/clients/gorm"
	redis_client "gin-realword-example/internal/modules/clients/redis"
	"gin-realword-example/internal/modules/core"
	"gin-realword-example/internal/modules/shared"
	"gin-realword-example/internal/modules/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var loginSessionExpireIn time.Duration

func init() {
	loginSessionExpireIn = core.ConfigStore.GetDuration(shared.ConfigKeyAuthLoginSessionExpireIn)
	err := gorm_client.GetDB().AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}
}

func UpsertUser(ctx context.Context, request models.CreateUserRequest) (*uint, error) {
	_, err := GetUserByEmail(ctx, request.Email)
	if err != nil {
		if err == ErrStoreUserNotFoundByEmail {
			return CreateUser(ctx, request)
		}
		return nil, err
	}
	return nil, nil
}

func CreateUser(ctx context.Context, request models.CreateUserRequest) (*uint, error) {
	user := models.User{
		Username: request.Username,
		Email:    request.Email,
	}
	err := gorm_client.GetDB().Create(&user).Error
	if err != nil {
		return nil, ErrStoreCreateUserFailed
	}
	return &user.ID, nil
}

func CheckUserExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := gorm_client.GetDB().Table("users").Where("id = ?", id).Count(&count).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return count > 0, nil
}

func GetUserByID(ctx context.Context, id uint) (models.User, error) {
	var user models.User
	err := gorm_client.GetDB().First(&user, id).Error
	if err != nil {
		return models.User{}, ErrStoreUserNotFoundByID
	}
	return user, nil
}

func GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := gorm_client.GetDB().First(&user, "email = ?", email).Error
	if err != nil {
		return models.User{}, ErrStoreUserNotFoundByEmail
	}
	return user, nil
}

func GenerateLoginSession(ctx context.Context, id uint) (*string, *time.Time, error) {
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
		SetEx(ctx, fmt.Sprintf(AuthKeyFormatLoginSession, sessionID), id, loginSessionExpireIn).
		Err()
	if err != nil {
		return nil, nil, err
	}
	return &sessionID, &expireAt, nil
}

func GetUserIDFromLoginSession(ctx context.Context, sessionID string) (*uint, *time.Time, error) {
	sessionKey := fmt.Sprintf(AuthKeyFormatLoginSession, sessionID)
	result, err := redis_client.GetRDB().Get(ctx, sessionKey).Result()
	if err != nil {
		return nil, nil, ErrStoreLoginSessionNotFound
	}
	id, err := strconv.Atoi(result)
	if err != nil {
		return nil, nil, err
	}
	expireAt := time.Now().Add(loginSessionExpireIn)
	err = redis_client.GetRDB().
		SetEx(ctx, sessionKey, id, loginSessionExpireIn).
		Err()
	if err != nil {
		return nil, nil, err
	}
	return utils.ToPtr(uint(id)), &expireAt, nil
}
