package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/usecase"
	"github.com/hiennguyen9874/stockk-go/internal/users"
	"github.com/hiennguyen9874/stockk-go/pkg/emailTemplates"
	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
	"github.com/hiennguyen9874/stockk-go/pkg/jwt"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/hiennguyen9874/stockk-go/pkg/secureRandom"
	"github.com/hiennguyen9874/stockk-go/pkg/sendEmail"
)

type userUseCase struct {
	usecase.UseCase[models.User]
	userPgRepo             users.UserPgRepository
	userRedisRepo          users.UserRedisRepository
	emailSender            sendEmail.EmailSender
	emailTemplateGenerator emailTemplates.EmailTemplatesGenerator
}

func CreateUserUseCaseI(
	userPgRepo users.UserPgRepository,
	userRedisRepo users.UserRedisRepository,
	cfg *config.Config,
	logger logger.Logger,
) users.UserUseCaseI {
	return &userUseCase{
		UseCase:                usecase.CreateUseCase[models.User](userPgRepo, cfg, logger),
		userPgRepo:             userPgRepo,
		userRedisRepo:          userRedisRepo,
		emailSender:            sendEmail.NewEmailSender(cfg),
		emailTemplateGenerator: emailTemplates.NewEmailTemplatesGenerator(cfg),
	}
}

func (u *userUseCase) Get(ctx context.Context, id uint) (*models.User, error) {
	cachedUser, err := u.userRedisRepo.GetObj(ctx, u.GenerateRedisUserKey(id))
	if err != nil {
		return nil, err
	}

	if cachedUser != nil {
		return cachedUser, nil
	}

	user, err := u.userPgRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = u.userRedisRepo.CreateObj(ctx, u.GenerateRedisUserKey(id), user, 3600); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) Delete(ctx context.Context, id uint) (*models.User, error) {
	user, err := u.userPgRepo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = u.userRedisRepo.Delete(ctx, u.GenerateRedisUserKey(id)); err != nil {
		return nil, err
	}

	if err = u.userRedisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id)); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) Update(
	ctx context.Context,
	id uint,
	values map[string]interface{},
) (*models.User, error) {
	obj, err := u.Get(ctx, id)
	if err != nil || obj == nil {
		return nil, err
	}

	user, err := u.userPgRepo.Update(ctx, obj, values)
	if err != nil {
		return nil, err
	}

	if err = u.userRedisRepo.Delete(ctx, u.GenerateRedisUserKey(id)); err != nil {
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

	user, err := u.userPgRepo.Create(ctx, exp)
	if err != nil {
		return nil, err
	}

	if user.Verified {
		return user, nil
	}

	verificationCode, err := secureRandom.RandomHex(16)
	if err != nil {
		return nil, err
	}

	// Update user in database
	updatedUser, err := u.userPgRepo.UpdateVerificationCode(ctx, user, verificationCode)
	if err != nil {
		return nil, err
	}

	bodyHtml, bodyPlain, err := u.emailTemplateGenerator.GenerateVerificationCodeTemplate(
		ctx,
		updatedUser.Name,
		fmt.Sprintf(
			"http://localhost:5000/auth/verifyemail?code=%s",
			verificationCode,
		),
	)
	if err != nil {
		return nil, err
	}

	err = u.emailSender.SendEmail(
		ctx,
		u.Cfg.Email.EmailFrom,
		updatedUser.Email,
		u.Cfg.Email.EmailVerificationSubject,
		bodyHtml,
		bodyPlain,
	)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (u *userUseCase) CreateUser(ctx context.Context, exp *models.User, confirmPassword string) (*models.User, error) {
	if exp.Password != confirmPassword {
		return nil, httpErrors.ErrValidation(errors.New("password do not match"))
	}
	return u.Create(ctx, exp)
}

func (u *userUseCase) createToken(ctx context.Context, exp models.User) (string, string, error) {
	accessToken, err := jwt.CreateAccessTokenRS256(
		exp.Id,
		exp.Email,
		u.Cfg.Jwt.JwtAccessTokenPrivateKey,
		u.Cfg.Jwt.JwtAccessTokenExpireDuration*int64(time.Minute),
		u.Cfg.Jwt.JwtIssuer,
	)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.CreateAccessTokenRS256(
		exp.Id,
		exp.Email,
		u.Cfg.Jwt.JwtRefreshTokenPrivateKey,
		u.Cfg.Jwt.JwtRefreshTokenExpireDuration*int64(time.Minute),
		u.Cfg.Jwt.JwtIssuer,
	)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *userUseCase) SignIn(ctx context.Context, email string, password string) (string, string, error) {
	user, err := u.userPgRepo.GetByEmail(ctx, email)
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

	if err = u.userRedisRepo.Sadd(
		ctx,
		u.GenerateRedisRefreshTokenKey(user.Id),
		refreshToken,
	); err != nil {
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
	user, err := u.userPgRepo.GetByEmail(ctx, u.Cfg.FirstSuperUser.FirstSuperUserEmail)

	if err != nil || user == nil {
		_, err := u.Create(ctx, &models.User{
			Name:        u.Cfg.FirstSuperUser.FirstSuperUserName,
			Email:       u.Cfg.FirstSuperUser.FirstSuperUserEmail,
			Password:    u.Cfg.FirstSuperUser.FirstSuperUserPassword,
			IsActive:    true,
			IsSuperUser: true,
			Verified:    true,
		})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (u *userUseCase) UpdatePassword(
	ctx context.Context,
	id uint,
	oldPassword string,
	newPassword string,
	confirmPassword string,
) (*models.User, error) {
	if newPassword != confirmPassword {
		return nil, httpErrors.ErrValidation(errors.New("password do not match"))
	}

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

	updatedUser, err := u.userPgRepo.UpdatePassword(ctx, user, hashedPassword)
	if err != nil {
		return nil, err
	}

	if err = u.userRedisRepo.Delete(ctx, u.GenerateRedisUserKey(id)); err != nil {
		return nil, err
	}

	if err = u.userRedisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id)); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (u *userUseCase) ParseIdFromRefreshToken(
	ctx context.Context,
	refreshToken string,
) (uint, error) {
	id, _, err := jwt.ParseTokenRS256(refreshToken, u.Cfg.Jwt.JwtRefreshTokenPublicKey)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *userUseCase) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	idParsed, err := u.ParseIdFromRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	isMember, err := u.userRedisRepo.SIsMember(
		ctx,
		u.GenerateRedisRefreshTokenKey(idParsed),
		refreshToken,
	)
	if err != nil {
		return "", "", err
	}

	if !isMember {
		return "", "",
			httpErrors.ErrNotFoundRefreshTokenRedis(errors.New("not found refresh token in redis"))
	}

	if err = u.userRedisRepo.Srem(
		ctx,
		u.GenerateRedisRefreshTokenKey(idParsed),
		refreshToken,
	); err != nil {
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

	if err = u.userRedisRepo.Sadd(
		ctx,
		u.GenerateRedisRefreshTokenKey(user.Id),
		refreshToken,
	); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err
}

func (u *userUseCase) Logout(ctx context.Context, refreshToken string) error {
	idParsed, err := u.ParseIdFromRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	if err = u.userRedisRepo.Srem(
		ctx,
		u.GenerateRedisRefreshTokenKey(idParsed),
		refreshToken,
	); err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) LogoutAll(ctx context.Context, id uint) error {
	if err := u.userRedisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id)); err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) Verify(ctx context.Context, verificationCode string) error {
	user, err := u.userPgRepo.GetByVerificationCode(ctx, verificationCode)
	if err != nil {
		return err
	}

	if user.Verified {
		return httpErrors.ErrUserAlreadyVerified(errors.New("user already verified"))
	}

	updatedUser, err := u.userPgRepo.UpdateVerification(ctx, user, "", true)
	if err != nil {
		return err
	}

	if err = u.userRedisRepo.Delete(ctx, u.GenerateRedisUserKey(updatedUser.Id)); err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) ForgotPassword(ctx context.Context, email string) error {
	user, err := u.userPgRepo.GetByEmail(ctx, email)

	if err != nil {
		return httpErrors.ErrNotFound(err)
	}

	if !user.Verified {
		return httpErrors.ErrUserNotVerified(errors.New("user not verified"))
	}

	resetToken, err := secureRandom.RandomHex(16)
	if err != nil {
		return err
	}

	updatedUser, err := u.userPgRepo.UpdatePasswordReset(
		ctx,
		user,
		resetToken,
		time.Now().Add(time.Minute*15),
	)
	if err != nil {
		return err
	}
	if err = u.userRedisRepo.Delete(ctx, u.GenerateRedisUserKey(updatedUser.Id)); err != nil {
		return err
	}

	bodyHtml, bodyPlain, err := u.emailTemplateGenerator.GeneratePasswordResetTemplate(
		ctx,
		updatedUser.Name,
		fmt.Sprintf("http://localhost:5000/auth/resetpassword?code=%s", resetToken),
	)
	if err != nil {
		return err
	}

	err = u.emailSender.SendEmail(
		ctx,
		u.Cfg.Email.EmailFrom,
		updatedUser.Email,
		u.Cfg.Email.EmailResetSubject,
		bodyHtml,
		bodyPlain,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) ResetPassword(
	ctx context.Context,
	resetToken string,
	newPassword string,
	confirmPassword string,
) error {
	if newPassword != confirmPassword {
		return httpErrors.ErrValidation(errors.New("password do not match"))
	}

	user, err := u.userPgRepo.GetByResetTokenResetAt(ctx, resetToken, time.Now())
	if err != nil {
		return err
	}

	hashedPassword, err := jwt.HashPassword(newPassword)
	if err != nil {
		return err
	}

	updatedUser, err := u.userPgRepo.UpdatePasswordResetToken(
		ctx,
		user,
		hashedPassword,
		"",
	)
	if err != nil {
		return err
	}

	if err = u.userRedisRepo.Delete(ctx, u.GenerateRedisUserKey(updatedUser.Id)); err != nil {
		return err
	}

	if err = u.userRedisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(updatedUser.Id)); err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) GenerateRedisUserKey(id uint) string {
	return fmt.Sprintf("%v:%v", models.User{}.TableName(), id)
}

func (u *userUseCase) GenerateRedisRefreshTokenKey(id uint) string {
	return fmt.Sprintf("RefreshToken:%v", id)
}
