package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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
	pgRepo    users.UserPgRepository
	redisRepo users.UserRedisRepository
}

func CreateUserUseCaseI(pgRepo users.UserPgRepository, redisRepo users.UserRedisRepository, cfg *config.Config, logger logger.Logger) users.UserUseCaseI {
	return &userUseCase{
		UseCase:   usecase.CreateUseCase[models.User](pgRepo, cfg, logger),
		pgRepo:    pgRepo,
		redisRepo: redisRepo,
	}
}

func (u *userUseCase) Get(ctx context.Context, id uuid.UUID) (*models.User, error) {
	cachedUser, err := u.redisRepo.Get(ctx, u.GenerateRedisUserKey(id))

	if err != nil {
		return nil, err
	}

	if cachedUser != nil {
		return cachedUser, nil
	}

	user, err := u.pgRepo.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.Create(ctx, u.GenerateRedisUserKey(id), user, 3600); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) Delete(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := u.pgRepo.Delete(ctx, id)

	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.Delete(ctx, u.GenerateRedisUserKey(id)); err != nil {
		return nil, err
	}

	if err = u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id)); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) Update(ctx context.Context, id uuid.UUID, values map[string]interface{}) (res *models.User, err error) {
	obj, err := u.Get(ctx, id)
	if err != nil || obj == nil {
		return nil, err
	}

	user, err := u.pgRepo.Update(ctx, obj, values)

	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.Delete(ctx, u.GenerateRedisUserKey(id)); err != nil {
		return nil, err
	}
	return user, nil
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

func (u *userUseCase) createToken(ctx context.Context, exp models.User) (string, string, error) {
	accessToken, err := jwt.CreateAccessTokenRS256(exp.Id.String(), exp.Email, u.Cfg.Jwt.JwtAccessTokenPrivateKey, u.Cfg.Jwt.JwtAccessTokenExpireDuration*int64(time.Minute), u.Cfg.Jwt.JwtIssuer)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.CreateAccessTokenRS256(exp.Id.String(), exp.Email, u.Cfg.Jwt.JwtRefreshTokenPrivateKey, u.Cfg.Jwt.JwtRefreshTokenExpireDuration*int64(time.Minute), u.Cfg.Jwt.JwtIssuer)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *userUseCase) SignIn(ctx context.Context, email string, password string) (string, string, error) {
	user, err := u.pgRepo.GetByEmail(ctx, email)

	if err != nil {
		return "", "", httpErrors.ErrNotFound(err)
	}

	if !jwt.ComparePassword(password, user.Password) {
		return "", "", httpErrors.ErrWrongPassword(errors.New("wrong password"))
	}

	accessToken, refreshToken, err := u.createToken(ctx, *user)
	if err != nil {
		return "", "", err
	}

	if err = u.redisRepo.Sadd(ctx, u.GenerateRedisRefreshTokenKey(user.Id), refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *userUseCase) IsActive(ctx context.Context, exp models.User) bool {
	return exp.IsActive
}

func (u *userUseCase) IsSuper(ctx context.Context, exp models.User) bool {
	return exp.IsSuperUser
}

func (u *userUseCase) CreateSuperUserIfNotExist(ctx context.Context) (bool, error) {
	user, err := u.pgRepo.GetByEmail(ctx, u.Cfg.FirstSuperUser.FirstSuperUserEmail)

	if err != nil || user == nil {
		_, err := u.Create(ctx, &models.User{
			Name:        u.Cfg.FirstSuperUser.FirstSuperUserName,
			Email:       u.Cfg.FirstSuperUser.FirstSuperUserEmail,
			Password:    u.Cfg.FirstSuperUser.FirstSuperUserPassword,
			IsActive:    true,
			IsSuperUser: true,
		})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (u *userUseCase) UpdatePassword(ctx context.Context, id uuid.UUID, oldPassword string, newPassword string) (*models.User, error) {
	user, err := u.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !jwt.ComparePassword(oldPassword, user.Password) {
		return nil, httpErrors.ErrWrongPassword(errors.New("old password and new password not same"))
	}

	hashedPassword, err := jwt.HashPassword(newPassword)

	if err != nil {
		return nil, err
	}

	updatedUser, err := u.pgRepo.UpdatePassword(ctx, user, hashedPassword)

	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.Delete(ctx, u.GenerateRedisUserKey(id)); err != nil {
		return nil, err
	}

	if err = u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id)); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (u *userUseCase) ParseIdFromRefreshToken(ctx context.Context, refreshToken string) (idParsed uuid.UUID, err error) {
	id, _, err := jwt.ParseTokenRS256(refreshToken, u.Cfg.Jwt.JwtRefreshTokenPublicKey)

	if err != nil {
		return uuid.UUID{}, err
	}

	idParsed, err = uuid.Parse(id)

	if err != nil {
		return uuid.UUID{}, httpErrors.ErrInvalidJWTClaims(errors.New("can not convert id to uuid from id in token"))
	}

	return idParsed, nil
}

func (u *userUseCase) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	idParsed, err := u.ParseIdFromRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	isMember, err := u.redisRepo.SIsMember(ctx, u.GenerateRedisRefreshTokenKey(idParsed), refreshToken)
	if err != nil {
		return "", "", err
	}

	if !isMember {
		return "", "", httpErrors.ErrNotFoundRefreshTokenRedis(errors.New("not found refresh token in redis"))
	}

	if err = u.redisRepo.Srem(ctx, u.GenerateRedisRefreshTokenKey(idParsed), refreshToken); err != nil {
		return "", "", err
	}

	user, err := u.Get(ctx, idParsed)

	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := u.createToken(ctx, *user)
	if err != nil {
		return "", "", err
	}

	if err = u.redisRepo.Sadd(ctx, u.GenerateRedisRefreshTokenKey(user.Id), refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err
}

func (u *userUseCase) Logout(ctx context.Context, refreshToken string) error {
	idParsed, err := u.ParseIdFromRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	if err = u.redisRepo.Srem(ctx, u.GenerateRedisRefreshTokenKey(idParsed), refreshToken); err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) LogoutAll(ctx context.Context, id uuid.UUID) error {
	if err := u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id)); err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) GenerateRedisUserKey(id uuid.UUID) string {
	return fmt.Sprintf("%s:%s", models.User{}.TableName(), id.String())
}

func (u *userUseCase) GenerateRedisRefreshTokenKey(id uuid.UUID) string {
	return fmt.Sprintf("RefreshToken:%s", id.String())
}
