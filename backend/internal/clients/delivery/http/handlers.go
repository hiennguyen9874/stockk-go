package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/clients"
	"github.com/hiennguyen9874/stockk-go/internal/clients/presenter"
	"github.com/hiennguyen9874/stockk-go/internal/middleware"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/hiennguyen9874/stockk-go/pkg/responses"
	"github.com/hiennguyen9874/stockk-go/pkg/utils"
)

type clientHandler struct {
	cfg       *config.Config
	clientsUC clients.ClientUseCaseI
	logger    logger.Logger
}

func CreateClientHandler(uc clients.ClientUseCaseI, cfg *config.Config, logger logger.Logger) clients.Handlers {
	return &clientHandler{cfg: cfg, clientsUC: uc, logger: logger}
}

// Create godoc
// @Summary Create Client
// @Description Create new client.
// @Tags clients
// @Accept json
// @Produce json
// @Param client body presenter.ClientCreate true "Add client"
// @Success 200 {object} responses.SuccessResponse[presenter.ClientResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /client [post]
func (h *clientHandler) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		client := new(presenter.ClientCreate)

		err := json.NewDecoder(r.Body).Decode(&client)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(ctx, client)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		newClient, err := h.clientsUC.CreateWithOwner(
			ctx,
			user.Id,
			mapModel(client),
		)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		clientResponse := *mapModelResponse(newClient)
		render.Respond(w, r, responses.CreateSuccessResponse(clientResponse))
	}
}

// Get godoc
// @Summary Read client
// @Description Get client by ID.
// @Tags clients
// @Accept json
// @Produce json
// @Param id path string true "Client Id"
// @Success 200 {object} responses.SuccessResponse[presenter.ClientResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /client/{id} [get]
func (h *clientHandler) Get() func(w http.ResponseWriter, r *http.Request) {
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

		client, err := h.clientsUC.Get(ctx, uint(id))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		if !user.IsSuperUser && client.OwnerId != user.Id {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(err))) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(client)))
	}
}

// GetMulti godoc
// @Summary Read Clients
// @Description Retrieve clients.
// @Tags clients
// @Accept json
// @Produce json
// @Param limit query int false "limit" Format(limit)
// @Param offset query int false "offset" Format(offset)
// @Success 200 {object} responses.SuccessResponse[[]presenter.ClientResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /client [get]
func (h *clientHandler) GetMulti() func(w http.ResponseWriter, r *http.Request) {
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

		var clients []*models.Client
		if user.IsSuperUser {
			clients, err = h.clientsUC.GetMulti(ctx, limit, offset)
		} else {
			clients, err = h.clientsUC.GetMultiByOwnerId(ctx, user.Id, limit, offset)
		}
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelsResponse(clients)))
	}
}

// Delete godoc
// @Summary Delete client
// @Description Delete an client by ID.
// @Tags clients
// @Accept json
// @Produce json
// @Param id path string true "Client Id"
// @Success 200 {object} responses.SuccessResponse[presenter.ClientResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /client/{id} [delete]
func (h *clientHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
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

		client, err := h.clientsUC.Get(ctx, uint(id))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		if !user.IsSuperUser && client.OwnerId != user.Id {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(err))) //nolint:errcheck
			return
		}

		err = h.clientsUC.DeleteWithoutGet(ctx, uint(id))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(client)))
	}
}

// Update godoc
// @Summary Update client
// @Description Update an client by ID.
// @Tags clients
// @Accept json
// @Produce json
// @Param id path string true "Client Id"
// @Param client body presenter.ClientUpdate true "Update client"
// @Success 200 {object} responses.SuccessResponse[presenter.ClientResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /client/{id} [put]
func (h *clientHandler) Update() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		client := new(presenter.ClientUpdate)

		err = json.NewDecoder(r.Body).Decode(&client)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		err = utils.ValidateStruct(r.Context(), client)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		dbClient, err := h.clientsUC.Get(ctx, uint(id))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		if !user.IsSuperUser && dbClient.OwnerId != user.Id {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(err))) //nolint:errcheck
			return
		}

		values := make(map[string]interface{})
		if client.CurrentTicker != nil && *client.CurrentTicker != "" {
			values["current_ticker"] = client.CurrentTicker
		}
		if client.CurrentResolution != nil && *client.CurrentResolution != "" {
			values["current_resolution"] = client.CurrentResolution
		}

		updatedClient, err := h.clientsUC.Update(r.Context(), uint(id), values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedClient)))
	}
}

func mapModel(exp *presenter.ClientCreate) *models.Client {
	return &models.Client{
		CurrentTicker:     exp.CurrentTicker,
		CurrentResolution: exp.CurrentResolution,
	}
}

func mapModelResponse(exp *models.Client) *presenter.ClientResponse {
	return &presenter.ClientResponse{
		Id:                exp.Id,
		CreatedAt:         exp.CreatedAt,
		UpdatedAt:         exp.UpdatedAt,
		CurrentTicker:     exp.CurrentTicker,
		CurrentResolution: exp.CurrentResolution,
		OwnerId:           exp.OwnerId,
	}
}

func mapModelsResponse(exp []*models.Client) []*presenter.ClientResponse {
	out := make([]*presenter.ClientResponse, len(exp))
	for i, user := range exp {
		out[i] = mapModelResponse(user)
	}
	return out
}
