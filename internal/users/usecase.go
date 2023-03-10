package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type UserUseCaseI interface {
	internal.UseCaseI[models.User]
	SignIn(ctx context.Context, email string, password string) (string, string, error)
	IsActive(ctx context.Context, exp models.User) bool
	IsSuper(ctx context.Context, exp models.User) bool
	CreateSuperUserIfNotExist(context.Context) (bool, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, oldPassword string, newPassword string) (*models.User, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
}
