package http

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/drawingtemplates"
	"github.com/hiennguyen9874/stockk-go/internal/drawingtemplates/presenter"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
)

type drawingTemplateHandler struct {
	cfg                *config.Config
	drawingTemplatesUC drawingtemplates.DrawingTemplateUseCaseI
	logger             logger.Logger
}

func CreateDrawingTemplateHandler(uc drawingtemplates.DrawingTemplateUseCaseI, cfg *config.Config, logger logger.Logger) drawingtemplates.Handlers {
	return &drawingTemplateHandler{cfg: cfg, drawingTemplatesUC: uc, logger: logger}
}

func (h *drawingTemplateHandler) CreateOrUpdate() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		q := r.URL.Query()
		clientId := q.Get("client")
		userId := q.Get("user")

		r.ParseMultipartForm(0) //nolint:errcheck
		templateName := r.FormValue("name")
		tool := r.FormValue("tool")
		content := r.FormValue("content")

		_, _, _, err := h.drawingTemplatesUC.CreateOrUpdateWithOwnerNameTool(ctx, clientId, userId, templateName, tool, content)
		if err != nil {
			render.Render(w, r, CreateDrawingTemplateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, presenter.DrawingTemplateCreateUpdateResponse{Status: "ok"})
	}
}

func (h *drawingTemplateHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		q := r.URL.Query()
		clientId := q.Get("client")
		userId := q.Get("user")
		templateName := q.Get("name")
		tool := q.Get("tool")

		if templateName == "" {
			studyTemplates, err := h.drawingTemplatesUC.GetAllByOwnerTool(ctx, clientId, userId, tool)
			if err != nil {
				render.Render(w, r, CreateDrawingTemplateErrorResponse(err)) //nolint:errcheck
				return
			}

			render.Respond(w, r, presenter.DrawingTemplateGetsResponse{
				Status: "ok",
				Data:   mapModelsResponse(studyTemplates),
			})
			return
		}

		studyTemplate, err := h.drawingTemplatesUC.GetByOwnerNameTool(r.Context(), clientId, userId, templateName, tool)
		if err != nil {
			render.Render(w, r, CreateDrawingTemplateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, presenter.DrawingTemplateGetResponse{
			Status: "ok",
			Data:   mapModelResponse(studyTemplate),
		})
	}
}

func (h *drawingTemplateHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		q := r.URL.Query()
		clientId := q.Get("client")
		userId := q.Get("user")
		templateName := q.Get("template")
		tool := q.Get("tool")

		_, err := h.drawingTemplatesUC.DeleteByOwnerNameTool(ctx, clientId, userId, templateName, tool)
		if err != nil {
			render.Render(w, r, CreateDrawingTemplateErrorResponse(err)) //nolint:errcheck
			return
		}

		render.Respond(w, r, presenter.DrawingTemplateDeleteResponse{Status: "ok"})
	}
}

func mapModelResponse(exp *models.DrawingTemplate) *presenter.DrawingTemplateResponse {
	return &presenter.DrawingTemplateResponse{
		Name:    exp.Name,
		Content: exp.Content,
	}
}

func mapModelsResponse(exp []*models.DrawingTemplate) []*string {
	out := make([]*string, len(exp))
	for i, studyTemplate := range exp {
		out[i] = &studyTemplate.Name
	}
	return out
}

func CreateDrawingTemplateErrorResponse(err error) render.Renderer {
	parsedErr := httpErrors.ParseErrors(err)

	return &presenter.DrawingTemplateErrorResponse{
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
