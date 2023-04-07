package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/charts"
	"github.com/hiennguyen9874/stockk-go/internal/charts/presenter"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/hiennguyen9874/stockk-go/pkg/utils"
)

type chartHandler struct {
	cfg      *config.Config
	chartsUC charts.ChartUseCaseI
	logger   logger.Logger
}

func CreateChartHandler(uc charts.ChartUseCaseI, cfg *config.Config, logger logger.Logger) charts.Handlers {
	return &chartHandler{cfg: cfg, chartsUC: uc, logger: logger}
}

func (h *chartHandler) CreateOrUpdate() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		q := r.URL.Query()
		chartId := q.Get("chart")
		clientId := q.Get("client")
		userId := q.Get("user")

		if chartId == "" {
			r.ParseMultipartForm(0) //nolint:errcheck

			chart := new(presenter.ChartCreate)
			chart.OwnerSource = clientId
			chart.OwnerId = userId
			chart.Name = r.FormValue("name")
			chart.Symbol = r.FormValue("symbol")
			chart.Resolution = r.FormValue("resolution")
			chart.Content = r.FormValue("content")

			err := utils.ValidateStruct(ctx, chart)
			if err != nil {
				render.Render(w, r, CreateChartErrorResponse(err)) //nolint:errcheck
				return
			}

			newChart, err := h.chartsUC.Create(ctx, mapModel(chart))
			if err != nil {
				render.Render(w, r, CreateChartErrorResponse(err)) //nolint:errcheck
				return
			}

			render.Respond(w, r, presenter.ChartCreateResponse{
				Status: "ok",
				Id:     newChart.Id,
			})
			return
		}

		id, err := strconv.ParseUint(chartId, 10, 32)
		if err != nil {
			render.Render(w, r, CreateChartErrorResponse(err)) //nolint:errcheck
			return
		}

		chart := new(presenter.ChartUpdate)
		chart.OwnerSource = clientId
		chart.OwnerId = userId
		chart.Name = r.FormValue("name")
		chart.Symbol = r.FormValue("symbol")
		chart.Resolution = r.FormValue("resolution")
		chart.Content = r.FormValue("content")

		values := make(map[string]interface{})
		if chart.OwnerSource != "" {
			values["owner_source"] = chart.OwnerSource
		}
		if chart.OwnerId != "" {
			values["owner_id"] = chart.OwnerId
		}
		if chart.Name != "" {
			values["name"] = chart.Name
		}
		if chart.Symbol != "" {
			values["symbol"] = chart.Symbol
		}
		if chart.Resolution != "" {
			values["resolution"] = chart.Resolution
		}
		if chart.Content != "" {
			values["content"] = chart.Content
		}

		_, err = h.chartsUC.Update(r.Context(), uint(id), values)
		if err != nil {
			render.Render(w, r, CreateChartErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, presenter.ChartUpdateResponse{Status: "ok"})
	}
}

func (h *chartHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		q := r.URL.Query()
		chartId := q.Get("chart")
		clientId := q.Get("client")
		userId := q.Get("user")

		if chartId == "" {
			charts, err := h.chartsUC.GetAllByOwner(ctx, clientId, userId)
			if err != nil {
				render.Render(w, r, CreateChartErrorResponse(err)) //nolint:errcheck
				return
			}

			render.Respond(w, r, presenter.ChartGetsResponse{
				Status: "ok",
				Data:   mapModelsResponse(charts),
			})
			return
		}

		id, err := strconv.ParseUint(chartId, 10, 32)
		if err != nil {
			render.Render(w, r, CreateChartErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		chart, err := h.chartsUC.Get(r.Context(), uint(id))
		if err != nil {
			render.Render(w, r, CreateChartErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, presenter.ChartGetResponse{
			Status: "ok",
			Data:   mapModelResponse(chart),
		})
	}
}

func (h *chartHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		q := r.URL.Query()
		chartId := q.Get("chart")
		clientId := q.Get("client")
		userId := q.Get("user")

		id, err := strconv.ParseUint(chartId, 10, 32)
		if err != nil {
			render.Render(w, r, CreateChartErrorResponse(httpErrors.ErrValidation(err))) //nolint:errcheck
			return
		}

		_, err = h.chartsUC.DeleteByIdOwner(ctx, uint(id), clientId, userId)
		if err != nil {
			render.Render(w, r, CreateChartErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, presenter.ChartDeleteResponse{Status: "ok"})
	}
}

func mapModel(exp *presenter.ChartCreate) *models.Chart {
	return &models.Chart{
		OwnerSource: exp.OwnerSource,
		OwnerId:     exp.OwnerId,
		Name:        exp.Name,
		Symbol:      exp.Symbol,
		Resolution:  exp.Resolution,
		Content:     exp.Content,
	}
}

func mapModelResponse(exp *models.Chart) *presenter.ChartResponse {
	return &presenter.ChartResponse{
		Id:        exp.Id,
		Name:      exp.Name,
		Timestamp: exp.LastModified.Unix(),
		Content:   exp.Content,
	}
}

func mapModelsResponse(exp []*models.Chart) []*presenter.ChartsResponse {
	out := make([]*presenter.ChartsResponse, len(exp))
	for i, chart := range exp {
		out[i] = &presenter.ChartsResponse{
			Id:         chart.Id,
			Name:       chart.Name,
			Symbol:     chart.Symbol,
			Resolution: chart.Resolution,
			Timestamp:  chart.LastModified.Unix(),
		}
	}
	return out
}

func CreateChartErrorResponse(err error) render.Renderer {
	parsedErr := httpErrors.ParseErrors(err)

	return &presenter.ChartErrorResponse{
		Error: &httpErrors.ErrResponse{
			Err:        parsedErr.GetErr(),
			Status:     parsedErr.GetStatus(),
			StatusText: parsedErr.GetStatusText(),
			Msg:        parsedErr.GetMsg(),
		},
		Status:  "error",
		Message: parsedErr.GetMsg(),
	}
}
