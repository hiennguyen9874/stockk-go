package users

import (
	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type UserUseCaseI interface {
	internal.UseCaseI[models.User]
}
