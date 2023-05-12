package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/stockssnapshot"
	"github.com/hiennguyen9874/stockk-go/internal/stockssnapshot/presenter"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/hiennguyen9874/stockk-go/pkg/responses"
)

type stockSnapshotHandler struct {
	cfg              *config.Config
	stocksSnapshotUC stockssnapshot.StockSnapshotUseCaseI
	logger           logger.Logger
}

func CreatestockSnapshotHandler(uc stockssnapshot.StockSnapshotUseCaseI, cfg *config.Config, logger logger.Logger) stockssnapshot.Handlers {
	return &stockSnapshotHandler{cfg: cfg, stocksSnapshotUC: uc, logger: logger}
}

// Get godoc
// @Summary Read stock snapshot by symbol
// @Description Get stock snapshot by symbol.
// @Tags stocksnapshot
// @Accept json
// @Produce json
// @Param symbol path string true "Ticker symbol"
// @Success 200 {object} responses.SuccessResponse[presenter.TickerResponse]
// @Failure 400	{object} responses.ErrorResponse
// @Failure 401	{object} responses.ErrorResponse
// @Failure 403	{object} responses.ErrorResponse
// @Failure 404	{object} responses.ErrorResponse
// @Failure 422	{object} responses.ErrorResponse
// @Security OAuth2Password
// @Router /stocksnapshot/{symbol} [get]
func (h *stockSnapshotHandler) GetStockSnapshotBySymbol() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		symbol := chi.URLParam(r, "symbol")

		stockSnapshot, err := h.stocksSnapshotUC.GetStockSnapshotBySymbol(r.Context(), symbol)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(stockSnapshot)))
	}
}

func mapModelResponse(exp *models.StockSnapshot) *presenter.StockSnapshotResponse {
	return &presenter.StockSnapshotResponse{
		Ticker:          exp.Ticker,
		BasicPrice:      exp.BasicPrice,
		CeilingPrice:    exp.CeilingPrice,
		FloorPrice:      exp.FloorPrice,
		AccumulatedVol:  exp.AccumulatedVol,
		AccumulatedVal:  exp.AccumulatedVal,
		MatchPrice:      exp.MatchPrice,
		MatchQtty:       exp.MatchQtty,
		HighestPrice:    exp.HighestPrice,
		LowestPrice:     exp.LowestPrice,
		BuyForeignQtty:  exp.BuyForeignQtty,
		SellForeignQtty: exp.SellForeignQtty,
		ProjectOpen:     exp.ProjectOpen,
		CurrentRoom:     exp.CurrentRoom,
		FloorCode:       exp.FloorCode,
		TotalRoom:       exp.TotalRoom,
	}
}

func mapModelsResponse(exp []*models.StockSnapshot) []*presenter.StockSnapshotResponse { //nolint:unused
	out := make([]*presenter.StockSnapshotResponse, len(exp))
	for i, ticker := range exp {
		out[i] = mapModelResponse(ticker)
	}
	return out
}
