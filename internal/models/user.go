package models

import (
	gwm_app "github.com/slhmy/go-webmods/app"
)

type CreateUserRequest struct {
	Username string `form:"username" json:"username" validate:"required"       example:"john_doe"`
	Email    string `form:"email"    json:"email"    validate:"required,email" example:"john@example.com"`
}

type CreateUserResponse struct {
	ID gwm_app.ID `form:"id" json:"id"`
}

type UpdateUserRequest struct {
	Username *string `form:"username" json:"username" example:"john_doe"`
	Email    *string `form:"email"    json:"email"    example:"john@example.com" validate:"email"`
}

type User struct {
	gwm_app.Model
	Username string `bson:"username" form:"username" json:"username" example:"john_doe"`
	Email    string `bson:"email"    form:"email"    json:"email"    example:"john@example.com"`
}
