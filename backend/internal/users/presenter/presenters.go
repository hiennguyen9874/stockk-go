package presenter

import (
	"time"
)

type UserCreate struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
}

type UserUpdate struct {
	Name string `json:"name"`
}

type UserUpdatePassword struct {
	OldPassword     string `json:"old_password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
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

type UserSignIn struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type Token struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
}

type PublicKey struct {
	PublicKeyAccessToken  string `json:"public_key_access_token,omitempty"`
	PublicKeyRefreshToken string `json:"public_key_refresh_token,omitempty"`
}

type ForgotPassword struct {
	Email string `json:"email" validate:"required"`
}

type ResetPassword struct {
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
}
