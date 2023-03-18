package users

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type UserUseCaseI interface {
	internal.UseCaseI[models.User]
	CreateUser(ctx context.Context, exp *models.User, confirmPassword string) (*models.User, error)
	SignIn(ctx context.Context, email string, password string) (string, string, error)
	IsActive(ctx context.Context, exp models.User) bool
	IsSuper(ctx context.Context, exp models.User) bool
	CreateSuperUserIfNotExist(context.Context) (bool, error)
	UpdatePassword(ctx context.Context, id uint, oldPassword string, newPassword string, confirmPassword string) (*models.User, error)
	ParseIdFromRefreshToken(ctx context.Context, refreshToken string) (uint, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	GenerateRedisUserKey(id uint) string
	GenerateRedisRefreshTokenKey(id uint) string
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, id uint) error
	Verify(ctx context.Context, verificationCode string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, resetToken string, newPassword string, confirmPassword string) error
}
