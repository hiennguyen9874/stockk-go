package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/render"
	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/bars"
	"github.com/hiennguyen9874/stockk-go/internal/dchart"
	"github.com/hiennguyen9874/stockk-go/internal/dchart/presenter"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/tickers"
	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/hiennguyen9874/stockk-go/pkg/responses"
	"github.com/hiennguyen9874/stockk-go/pkg/utils"
)

type dchartHandler struct {
	cfg       *config.Config
	tickersUC tickers.TickerUseCaseI
	barUC     bars.BarUseCaseI
	logger    logger.Logger
}

func CreateDchartHandler(uc tickers.TickerUseCaseI, cfg *config.Config, logger logger.Logger) dchart.Handlers {
	return &dchartHandler{cfg: cfg, tickersUC: uc, logger: logger}
}

func (h *dchartHandler) GetTime() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		loc, _ := time.LoadLocation(h.cfg.Server.TimeZone)

		render.Respond(w, r, time.Now().In(loc))
	}
}

func (h *dchartHandler) GetConfig() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Get from config
		render.Respond(w, r, presenter.DchartConfig{
			Exchanges: &[]presenter.DchartExchange{
				{
					Name:  "All Exchanges",
					Value: "",
					Desc:  "",
				},
				{
					Name:  "HOSE",
					Value: "HOSE",
					Desc:  "Ho Chi Minh Stock Exchange",
				},
				{
					Name:  "HNX",
					Value: "HNX",
					Desc:  "Hanoi Stock Exchange",
				},
				{
					Name:  "UPCOM",
					Value: "UPCOM",
					Desc:  "Unlisted Public Company Market",
				},
			},
			SupportedResolutions:   &[]string{"1", "5", "15", "30", "60", "D", "W", "M"},
			SupportsMarks:          utils.NewBool(false),
			SupportsTime:           utils.NewBool(true),
			SupportsTimescaleMarks: utils.NewBool(false),
			SymbolsTypes: &[]presenter.DchartSymbolType{
				{
					Name:  "All types",
					Value: "",
				},
				{
					Name:  "Stock",
					Value: "stock",
				},
				{
					Name:  "Index",
					Value: "index",
				},
				{
					Name:  "Crypto",
					Value: "crypto",
				},
			},
			SupportsSearch:       utils.NewBool(true),
			SupportsGroupRequest: utils.NewBool(false),
		})
	}
}

func (h *dchartHandler) GetSymbols() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		q := r.URL.Query()
		symbol := q.Get("symbol")

		ticker, err := h.tickersUC.GetBySymbol(ctx, symbol)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, presenter.DchartLibrarySymbolInfo{
			Name:                 ticker.ShortName,
			FullName:             ticker.FullName,
			Ticker:               &ticker.Symbol,
			Description:          ticker.FullName,
			Type:                 "stock",
			Session:              "0900-1130,1300-1500",
			Exchange:             ticker.Exchange,
			ListedExchange:       ticker.Exchange,
			Timezone:             "Asia/Ho_Chi_Minh",
			Format:               "price",
			Pricescale:           100,
			Minmov:               1,
			Minmove2:             0,
			HasIntraday:          utils.NewBool(true),
			SupportedResolutions: []string{"1", "5", "15", "30", "60", "D", "W", "M"},
			IntradayMultipliers:  &[]string{"1", "5", "15", "30", "60"},
			HasDaily:             utils.NewBool(true),
			HasWeeklyAndMonthly:  utils.NewBool(true),
			HasEmptyBars:         utils.NewBool(false),
			HasNoVolume:          utils.NewBool(false),
		})
	}
}

func (h *dchartHandler) Search() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		q := r.URL.Query()
		limitQ, err := strconv.Atoi(q.Get("limit"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		queryQ := q.Get("query")
		typeQ := q.Get("type")
		exchangeQ := q.Get("type")

		if typeQ != "stock" {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrBadRequest(fmt.Errorf("not support type: %v", typeQ))))
			return
		}

		tickers, err := h.tickersUC.SearchBySymbol(ctx, queryQ, limitQ, exchangeQ)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		resultItems := make([]presenter.DchartSearchSymbolResultItem, len(tickers))
		for i, ticker := range tickers {
			resultItems[i] = presenter.DchartSearchSymbolResultItem{
				Symbol:      ticker.Symbol,
				FullName:    ticker.FullName,
				Description: ticker.FullName,
				Exchange:    ticker.Exchange,
				Ticker:      ticker.Symbol,
				Type:        "stock",
			}
		}

		render.Respond(w, r, resultItems)
	}
}

func (h *dchartHandler) History() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		loc, _ := time.LoadLocation(h.cfg.Server.TimeZone)

		q := r.URL.Query()
		symbolQ := q.Get("symbol")

		resolutionQ := q.Get("resolution")
		if resolutionQ != "D" {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrBadRequest(fmt.Errorf("not support resolution: %v", resolutionQ))))
			return
		}

		fromQ, err := strconv.ParseInt(q.Get("from"), 10, 64)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		toQ, err := strconv.ParseInt(q.Get("to"), 10, 64)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		countbackQ, err := strconv.Atoi(q.Get("countback"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		var bars []*models.Bar
		if countbackQ != 0 {
			bars, err = h.barUC.GetByToLimit(ctx, resolutionQ, symbolQ, utils.UpdateTimeZone(time.Unix(toQ, 0), loc).UTC(), countbackQ)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
		} else {
			bars, err = h.barUC.GetByFromTo(ctx, resolutionQ, symbolQ, utils.UpdateTimeZone(time.Unix(fromQ, 0), loc).UTC(), utils.UpdateTimeZone(time.Unix(toQ, 0), loc).UTC())
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(err))
				return
			}
		}

		var history presenter.DchartHistoryFullDataResponse
		history.Status = "ok"
		history.Time = make([]int64, len(bars))
		history.Close = make([]float64, len(bars))
		history.Open = make([]float64, len(bars))
		history.High = make([]float64, len(bars))
		history.Low = make([]float64, len(bars))
		history.Volume = make([]int64, len(bars))

		for i, bar := range bars {
			history.Time[i] = bar.Time.In(loc).Unix()
			history.Close[i] = bar.Close
			history.Open[i] = bar.Open
			history.High[i] = bar.High
			history.Low[i] = bar.Low
			history.Volume[i] = bar.Volume
		}
		render.Respond(w, r, history)
	}
}
