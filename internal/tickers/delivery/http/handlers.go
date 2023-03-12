package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/hiennguyen9874/go-boilerplate/config"
	"github.com/hiennguyen9874/go-boilerplate/internal/models"
	"github.com/hiennguyen9874/go-boilerplate/internal/tickers"
	"github.com/hiennguyen9874/go-boilerplate/internal/tickers/presenter"
	"github.com/hiennguyen9874/go-boilerplate/pkg/httpErrors"
	"github.com/hiennguyen9874/go-boilerplate/pkg/logger"
	"github.com/hiennguyen9874/go-boilerplate/pkg/responses"
)

type tickerHandler struct {
	cfg       *config.Config
	tickersUC tickers.TickerUseCaseI
	logger    logger.Logger
}

func CreateTickerHandler(uc tickers.TickerUseCaseI, cfg *config.Config, logger logger.Logger) tickers.Handlers {
	return &tickerHandler{cfg: cfg, tickersUC: uc, logger: logger}
}

func (h *tickerHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		ticker, err := h.tickersUC.Get(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(ticker)))
	}
}

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

func (h *tickerHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		ticker, err := h.tickersUC.Delete(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(ticker)))
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
	}
}

func mapModelsResponse(exp []*models.Ticker) []*presenter.TickerResponse {
	out := make([]*presenter.TickerResponse, len(exp))
	for i, user := range exp {
		out[i] = mapModelResponse(user)
	}
	return out
}
