package presenter

import (
	"time"
)

type UserCreate struct {
	Name            string `json:"name" validate:"required" example:"Xuan Hien"`
	Email           string `json:"email" validate:"required" example:"hiennguyen9874@gmail.com"`
	Password        string `json:"password" validate:"required,min=8" example:"password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8" example:"password"`
}

type UserUpdate struct {
	Name string `json:"name" example:"Xuan Hien"`
}

type UserUpdatePassword struct {
	OldPassword     string `json:"old_password" validate:"required,min=8" example:"old_password"`
	NewPassword     string `json:"new_password" validate:"required,min=8" example:"password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8" example:"password"`
}

type UserResponse struct {
	Id          uint      `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Email       string    `json:"email,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsActive    bool      `json:"is_active"`
	IsSuperUser bool      `json:"is_superuser"`
	Verified    bool      `json:"verified"`
}
