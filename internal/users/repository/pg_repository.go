package repository

import (
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/repository"
	"github.com/hiennguyen9874/stockk-go/internal/users"
	"gorm.io/gorm"
)

type UserRepo struct {
	repository.PgRepo[models.User]
}

func CreateUserRepository(db *gorm.DB) users.UserRepository {
	return &UserRepo{
		PgRepo: repository.CreatePgRepo[models.User](db),
	}
}
