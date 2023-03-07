package usecase

import (
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/usecase"
	"github.com/hiennguyen9874/stockk-go/internal/users"
)

type userUseCase struct {
	usecase.UseCase[models.User]
}

func CreateUserUseCaseI(repo users.UserRepository) users.UserUseCaseI {
	return &userUseCase{
		UseCase: usecase.CreateUseCase[models.User](repo),
	}
}
