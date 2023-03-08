package users

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type UserUseCaseI interface {
	internal.UseCaseI[models.User]
	SignIn(ctx context.Context, email string, password string) (string, error)
}
