package users

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type UserRepository interface {
	internal.PgRepository[models.User]
	GetByEmail(ctx context.Context, email string) (res *models.User, err error)
}
