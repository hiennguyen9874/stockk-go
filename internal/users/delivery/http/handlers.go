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
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
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
			render.Render(w, r, httpErrors.ErrRender(httpErrors.ParseErrors(err)))
			return
		}

		err = utils.ValidateStruct(r.Context(), user)

		if err != nil {
			render.Render(w, r, httpErrors.ErrRender(httpErrors.ErrValidation(err)))
			return
		}

		newUser, err := h.usersUC.Create(r.Context(), mapModel(user))
		if err != nil {
			render.Render(w, r, httpErrors.ErrRender(httpErrors.ParseErrors(err)))
			return
		}
		render.Respond(w, r, mapModelResponse(newUser))
	}
}

func (h *userHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))

		if err != nil {
			render.Render(w, r, httpErrors.ErrRender(httpErrors.ErrValidation(err)))
			return
		}
		user, err := h.usersUC.Get(r.Context(), id)
		if err != nil {
			render.Render(w, r, httpErrors.ErrRender(httpErrors.ParseErrors(err)))
			return
		}
		render.Respond(w, r, mapModelResponse(user))
	}
}

func (h *userHandler) GetMulti() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		limit, _ := strconv.Atoi(q.Get("limit"))
		offset, _ := strconv.Atoi(q.Get("offset"))

		users, err := h.usersUC.GetMulti(r.Context(), limit, offset)

		if err != nil {
			render.Render(w, r, httpErrors.ErrRender(httpErrors.ParseErrors(err)))
			return
		}

		render.Respond(w, r, mapModelsResponse(users))
	}
}

func (h *userHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))

		if err != nil {
			render.Render(w, r, httpErrors.ErrRender(httpErrors.ErrValidation(err)))
			return
		}
		user, err := h.usersUC.Delete(r.Context(), id)
		if err != nil {
			render.Render(w, r, httpErrors.ErrRender(httpErrors.ParseErrors(err)))
			return
		}
		render.Respond(w, r, mapModelResponse(user))
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
			render.Render(w, r, httpErrors.ErrRender(httpErrors.ErrValidation(err)))
			return
		}

		token, err := h.usersUC.SignIn(r.Context(), user.Email, user.Password)

		if err != nil {
			render.Render(w, r, httpErrors.ErrRender(httpErrors.ParseErrors(err)))
			return
		}

		render.Respond(w, r, presenter.Token{AccessToken: token, TokenType: "bearer"})
	}
}

func (h *userHandler) Me() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := middleware.GetUserFromCtx(ctx)

		if err != nil {
			render.Render(w, r, httpErrors.ErrRender(httpErrors.ParseErrors(err)))
			return
		}

		render.Respond(w, r, mapModelResponse(user))
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
