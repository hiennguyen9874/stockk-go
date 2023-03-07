package users

import (
	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type UserRepository interface {
	internal.PgRepository[models.User]
}
