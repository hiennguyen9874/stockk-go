package users

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type UserPgRepository interface {
	internal.PgRepository[models.User]
	GetByEmail(ctx context.Context, email string) (res *models.User, err error)
	UpdatePassword(ctx context.Context, exp *models.User, newPassword string) (res *models.User, err error)
}
