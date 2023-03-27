package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/tickers"
	"github.com/hiennguyen9874/stockk-go/internal/tickers/presenter"
	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/hiennguyen9874/stockk-go/pkg/responses"
)

type tickerHandler struct {
	cfg       *config.Config
	tickersUC tickers.TickerUseCaseI
	logger    logger.Logger
}

func CreateTickerHandler(uc tickers.TickerUseCaseI, cfg *config.Config, logger logger.Logger) tickers.Handlers {
	return &tickerHandler{cfg: cfg, tickersUC: uc, logger: logger}
}

// GetMulti godoc
// @Summary Read tickers
// @Description Retrieve tickers.
// @Tags tickers
// @Accept json
// @Produce json
// @Param limit query int false "limit" Format(limit)
// @Param offset query int false "offset" Format(offset)
// @Success 200 {object} responses.Response
// @Failure 400	{object} responses.Response
// @Failure 401	{object} responses.Response
// @Failure 422	{object} responses.Response
// @Security OAuth2Password
// @Router /ticker [get]
func (h *tickerHandler) GetMulti() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		limit, _ := strconv.Atoi(q.Get("limit"))
		offset, _ := strconv.Atoi(q.Get("offset"))

		tickers, err := h.tickersUC.GetMulti(r.Context(), limit, offset)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelsResponse(tickers)))
	}
}

// Get godoc
// @Summary Read ticker
// @Description Get ticker by symbol.
// @Tags tickers
// @Accept json
// @Produce json
// @Param symbol path string true "Ticker symbol"
// @Success 200 {object} responses.Response
// @Failure 400	{object} responses.Response
// @Failure 401	{object} responses.Response
// @Failure 403	{object} responses.Response
// @Failure 404	{object} responses.Response
// @Failure 422	{object} responses.Response
// @Security OAuth2Password
// @Router /ticker/{symbol} [get]
func (h *tickerHandler) GetBySymbol() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		symbol := chi.URLParam(r, "symbol")

		ticker, err := h.tickersUC.GetBySymbol(r.Context(), symbol)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(ticker)))
	}
}

// Update godoc
// @Summary Update ticker
// @Description Update an ticker by Symbol.
// @Tags tickers
// @Accept json
// @Produce json
// @Param symbol path string true "Ticker symbol"
// @Param is_active query bool false "is_active" Format(is_active)
// @Success 200 {object} responses.Response
// @Failure 400	{object} responses.Response
// @Failure 401	{object} responses.Response
// @Failure 403	{object} responses.Response
// @Failure 404	{object} responses.Response
// @Failure 422	{object} responses.Response
// @Security OAuth2Password
// @Router /ticker/{symbol} [put]
func (h *tickerHandler) UpdateIsActiveBySymbol() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		symbol := chi.URLParam(r, "symbol")

		q := r.URL.Query()
		isActiveString := q.Get("is_active")
		isActive, err := strconv.ParseBool(isActiveString)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		updatedTicker, err := h.tickersUC.UpdateIsActiveBySymbol(ctx, symbol, isActive)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedTicker)))
	}
}

func mapModelResponse(exp *models.Ticker) *presenter.TickerResponse {
	return &presenter.TickerResponse{
		Id:        exp.Id,
		Symbol:    exp.Symbol,
		Exchange:  exp.Exchange,
		FullName:  exp.FullName,
		ShortName: exp.ShortName,
		Type:      exp.Type,
		IsActive:  exp.IsActive,
	}
}

func mapModelsResponse(exp []*models.Ticker) []*presenter.TickerResponse {
	out := make([]*presenter.TickerResponse, len(exp))
	for i, ticker := range exp {
		out[i] = mapModelResponse(ticker)
	}
	return out
}
