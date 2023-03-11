package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/middleware"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/users"
	"github.com/hiennguyen9874/stockk-go/internal/users/presenter"
	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
	"github.com/hiennguyen9874/stockk-go/pkg/jwt"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/hiennguyen9874/stockk-go/pkg/responses"
	"github.com/hiennguyen9874/stockk-go/pkg/utils"
)

type userHandler struct {
	cfg     *config.Config
	usersUC users.UserUseCaseI
	logger  logger.Logger
}

func CreateUserHandler(uc users.UserUseCaseI, cfg *config.Config, logger logger.Logger) users.Handlers {
	return &userHandler{cfg: cfg, usersUC: uc, logger: logger}
}

func (h *userHandler) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := new(presenter.UserCreate)

		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		err = utils.ValidateStruct(r.Context(), user)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		newUser, err := h.usersUC.Create(r.Context(), mapModel(user))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		userResponse := *mapModelResponse(newUser)

		render.Respond(w, r, responses.CreateSuccessResponse(userResponse))
	}
}

func (h *userHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		user, err := h.usersUC.Get(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(user)))
	}
}

func (h *userHandler) GetMulti() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		limit, _ := strconv.Atoi(q.Get("limit"))
		offset, _ := strconv.Atoi(q.Get("offset"))

		users, err := h.usersUC.GetMulti(r.Context(), limit, offset)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelsResponse(users)))
	}
}

func (h *userHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		user, err := h.usersUC.Delete(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(user)))
	}
}

func (h *userHandler) Update() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		user := new(presenter.UserUpdate)

		err = json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		err = utils.ValidateStruct(r.Context(), user)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		values := make(map[string]interface{})
		if user.Name != "" {
			values["name"] = user.Name
		}

		updatedUser, err := h.usersUC.Update(r.Context(), id, values)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

func (h *userHandler) UpdatePassword() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		user := new(presenter.UserUpdatePassword)

		err = json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		err = utils.ValidateStruct(r.Context(), user)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		updatedUser, err := h.usersUC.UpdatePassword(r.Context(), id, user.OldPassword, user.NewPassword)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

func (h *userHandler) SignIn() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := new(presenter.UserSignIn)

		r.ParseMultipartForm(0)
		user.Email = r.FormValue("email")
		user.Password = r.FormValue("password")

		err := utils.ValidateStruct(r.Context(), user)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		accessToken, refreshToken, err := h.usersUC.SignIn(r.Context(), user.Email, user.Password)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, presenter.Token{AccessToken: accessToken, RefreshToken: refreshToken, TokenType: "bearer"})
	}
}

func (h *userHandler) Me() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := middleware.GetUserFromCtx(ctx)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(user)))
	}
}

func (h *userHandler) UpdateMe() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := middleware.GetUserFromCtx(ctx)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		userUpdate := new(presenter.UserUpdate)

		err = json.NewDecoder(r.Body).Decode(&userUpdate)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		err = utils.ValidateStruct(r.Context(), userUpdate)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		values := make(map[string]interface{})
		if userUpdate.Name != "" {
			values["name"] = userUpdate.Name
		}

		updatedUser, err := h.usersUC.Update(r.Context(), user.Id, values)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

func (h *userHandler) UpdatePasswordMe() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := middleware.GetUserFromCtx(ctx)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		userUpdate := new(presenter.UserUpdatePassword)

		err = json.NewDecoder(r.Body).Decode(&userUpdate)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		err = utils.ValidateStruct(r.Context(), userUpdate)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		updatedUser, err := h.usersUC.UpdatePassword(r.Context(), user.Id, userUpdate.OldPassword, userUpdate.NewPassword)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedUser)))
	}
}

func (h *userHandler) RefreshToken() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		refreshToken := middleware.TokenFromHeader(r)

		accessToken, refreshToken, err := h.usersUC.Refresh(ctx, refreshToken)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, presenter.Token{AccessToken: accessToken, RefreshToken: refreshToken, TokenType: "bearer"})
	}
}

func (h *userHandler) GetPublicKey() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		publicKeyAccessToken, err := jwt.DecodeBase64(h.cfg.Jwt.JwtAccessTokenPublicKey)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		publicKeyRefreshToken, err := jwt.DecodeBase64(h.cfg.Jwt.JwtRefreshTokenPublicKey)

		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, presenter.PublicKey{
			PublicKeyAccessToken:  string(publicKeyAccessToken[:]),
			PublicKeyRefreshToken: string(publicKeyRefreshToken[:]),
		})
	}
}

func mapModel(exp *presenter.UserCreate) *models.User {
	return &models.User{
		Name:        exp.Name,
		Email:       exp.Email,
		Password:    exp.Password,
		IsActive:    true,
		IsSuperUser: false,
	}
}

func mapModelResponse(exp *models.User) *presenter.UserResponse {
	return &presenter.UserResponse{
		Id:          exp.Id,
		Name:        exp.Name,
		Email:       exp.Email,
		CreatedAt:   exp.CreatedAt,
		UpdatedAt:   exp.UpdatedAt,
		IsActive:    exp.IsActive,
		IsSuperUser: exp.IsSuperUser,
	}
}

func mapModelsResponse(exp []*models.User) []*presenter.UserResponse {
	out := make([]*presenter.UserResponse, len(exp))
	for i, user := range exp {
		out[i] = mapModelResponse(user)
	}
	return out
}
