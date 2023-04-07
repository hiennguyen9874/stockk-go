package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/middleware"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/watchlists"
	"github.com/hiennguyen9874/stockk-go/internal/watchlists/presenter"
	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/hiennguyen9874/stockk-go/pkg/responses"
	"github.com/hiennguyen9874/stockk-go/pkg/utils"
	"github.com/lib/pq"
)

type watchListHandler struct {
	cfg          *config.Config
	watchListsUC watchlists.WatchListUseCaseI
	logger       logger.Logger
}

func CreateWatchListHandler(uc watchlists.WatchListUseCaseI, cfg *config.Config, logger logger.Logger) watchlists.Handlers {
	return &watchListHandler{cfg: cfg, watchListsUC: uc, logger: logger}
}

// Create godoc
// @Summary Create WatchList
// @Description Create new watchList.
// @Tags watchlists
// @Accept json
// @Produce json
// @Param watchList body presenter.WatchListCreate true "Add watchList"
// @Success 200 {object} responses.SuccessResponse[presenter.WatchListResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /watchlist [post]
func (h *watchListHandler) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		watchList := new(presenter.WatchListCreate)

		err := json.NewDecoder(r.Body).Decode(&watchList)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(ctx, watchList)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		newWatchList, err := h.watchListsUC.CreateWithOwner(
			ctx,
			user.Id,
			mapModel(watchList),
		)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		watchListResponse := *mapModelResponse(newWatchList)
		render.Respond(w, r, responses.CreateSuccessResponse(watchListResponse))
	}
}

// Get godoc
// @Summary Read watchList
// @Description Get watchList by ID.
// @Tags watchlists
// @Accept json
// @Produce json
// @Param id path string true "WatchList Id"
// @Success 200 {object} responses.SuccessResponse[presenter.WatchListResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /watchlist/{id} [get]
func (h *watchListHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		watchList, err := h.watchListsUC.Get(ctx, uint(id))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		if !user.IsSuperUser && watchList.OwnerId != user.Id {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(err))) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(watchList)))
	}
}

// GetMulti godoc
// @Summary Read WatchLists
// @Description Retrieve watchlists.
// @Tags watchlists
// @Accept json
// @Produce json
// @Param limit query int false "limit" Format(limit)
// @Param offset query int false "offset" Format(offset)
// @Success 200 {object} responses.SuccessResponse[[]presenter.WatchListResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /watchlist [get]
func (h *watchListHandler) GetMulti() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		limit, _ := strconv.Atoi(q.Get("limit"))
		offset, _ := strconv.Atoi(q.Get("offset"))

		ctx := r.Context()

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		var watchLists []*models.WatchList
		if user.IsSuperUser {
			watchLists, err = h.watchListsUC.GetMulti(ctx, limit, offset)
		} else {
			watchLists, err = h.watchListsUC.GetMultiByOwnerId(ctx, user.Id, limit, offset)
		}
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelsResponse(watchLists)))
	}
}

// Delete godoc
// @Summary Delete watchList
// @Description Delete an watchList by ID.
// @Tags watchlists
// @Accept json
// @Produce json
// @Param id path string true "WatchList Id"
// @Success 200 {object} responses.SuccessResponse[presenter.WatchListResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /watchlist/{id} [delete]
func (h *watchListHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		watchList, err := h.watchListsUC.Get(ctx, uint(id))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		if !user.IsSuperUser && watchList.OwnerId != user.Id {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(err))) //nolint:errcheck
			return
		}

		err = h.watchListsUC.DeleteWithoutGet(ctx, uint(id))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(watchList)))
	}
}

// Update godoc
// @Summary Update watchList
// @Description Update an watchList by ID.
// @Tags watchlists
// @Accept json
// @Produce json
// @Param id path string true "WatchList Id"
// @Param watchList body presenter.WatchListUpdate true "Update watchList"
// @Success 200 {object} responses.SuccessResponse[presenter.WatchListResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /watchlist/{id} [put]
func (h *watchListHandler) Update() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		watchList := new(presenter.WatchListUpdate)

		err = json.NewDecoder(r.Body).Decode(&watchList)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(r.Context(), watchList)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		dbWatchList, err := h.watchListsUC.Get(ctx, uint(id))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		if !user.IsSuperUser && dbWatchList.OwnerId != user.Id {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(err))) //nolint:errcheck
			return
		}

		values := make(map[string]interface{})
		if watchList.Name != nil {
			values["name"] = watchList.Name
		}
		if watchList.Tickers != nil {
			values["tickers"] = pq.StringArray(*watchList.Tickers)
		}

		updatedWatchList, err := h.watchListsUC.Update(r.Context(), uint(id), values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedWatchList)))
	}
}

func mapModel(exp *presenter.WatchListCreate) *models.WatchList {
	return &models.WatchList{
		Name:    exp.Name,
		Tickers: exp.Tickers,
	}
}

func mapModelResponse(exp *models.WatchList) *presenter.WatchListResponse {
	return &presenter.WatchListResponse{
		Id:        exp.Id,
		CreatedAt: exp.CreatedAt,
		UpdatedAt: exp.UpdatedAt,
		Name:      exp.Name,
		Tickers:   exp.Tickers,
		OwnerId:   exp.OwnerId,
	}
}

func mapModelsResponse(exp []*models.WatchList) []*presenter.WatchListResponse {
	out := make([]*presenter.WatchListResponse, len(exp))
	for i, user := range exp {
		out[i] = mapModelResponse(user)
	}
	return out
}
