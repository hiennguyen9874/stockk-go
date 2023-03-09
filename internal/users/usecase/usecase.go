package usecase

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/usecase"
	"github.com/hiennguyen9874/stockk-go/internal/users"
	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
	"github.com/hiennguyen9874/stockk-go/pkg/jwt"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
)

type userUseCase struct {
	usecase.UseCase[models.User]
	pgRepo users.UserRepository
}

func CreateUserUseCaseI(repo users.UserRepository, cfg *config.Config, logger logger.Logger) users.UserUseCaseI {
	return &userUseCase{
		UseCase: usecase.CreateUseCase[models.User](repo, cfg, logger),
		pgRepo:  repo,
	}
}

func (u *userUseCase) Create(ctx context.Context, exp *models.User) (*models.User, error) {
	exp.Email = strings.ToLower(strings.TrimSpace(exp.Email))
	exp.Password = strings.TrimSpace(exp.Password)

	hashedPassword, err := jwt.HashPassword(exp.Password)

	if err != nil {
		return nil, err
	}
	exp.Password = hashedPassword

	return u.pgRepo.Create(ctx, exp)
}

func (u *userUseCase) SignIn(ctx context.Context, email string, password string) (string, error) {
	user, err := u.pgRepo.GetByEmail(ctx, email)

	if err != nil {
		return "", httpErrors.ErrNotFound(err)
	}

	if !jwt.ComparePassword(password, user.Password) {
		return "", httpErrors.Err(httpErrors.ErrorWrongPassword, http.StatusBadRequest, httpErrors.ErrorWrongPassword.Error())
	}

	return jwt.CreateAccessToken(user.Id.String(), user.Email, u.Cfg.Server.JwtSecretKey, u.Cfg.Server.JwtExpireDuration*int64(time.Minute), u.Cfg.Server.JwtIssuer)
}

func (u *userUseCase) IsActive(ctx context.Context, exp models.User) bool {
	return exp.IsActive
}

func (u *userUseCase) IsSuper(ctx context.Context, exp models.User) bool {
	return exp.IsSuperUser
}
