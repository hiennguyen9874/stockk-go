package http

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/studytemplates"
	"github.com/hiennguyen9874/stockk-go/internal/studytemplates/presenter"
	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
)

type studyTemplateHandler struct {
	cfg              *config.Config
	studyTemplatesUC studytemplates.StudyTemplateUseCaseI
	logger           logger.Logger
}

func CreateStudyTemplateHandler(uc studytemplates.StudyTemplateUseCaseI, cfg *config.Config, logger logger.Logger) studytemplates.Handlers {
	return &studyTemplateHandler{cfg: cfg, studyTemplatesUC: uc, logger: logger}
}

func (h *studyTemplateHandler) CreateOrUpdate() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		q := r.URL.Query()
		clientId := q.Get("client")
		userId := q.Get("user")

		r.ParseMultipartForm(0)
		templateName := r.FormValue("name")
		content := r.FormValue("content")

		_, _, _, err := h.studyTemplatesUC.CreateOrUpdateWithOwnerName(ctx, clientId, userId, templateName, content)
		if err != nil {
			render.Render(w, r, CreateStudyTemplateErrorResponse(err))
			return
		}

		render.Respond(w, r, presenter.StudyTemplateCreateUpdateResponse{Status: "ok"})
	}
}

func (h *studyTemplateHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		q := r.URL.Query()
		clientId := q.Get("client")
		userId := q.Get("user")
		templateName := q.Get("template")

		if templateName == "" {
			studyTemplates, err := h.studyTemplatesUC.GetAllByOwner(ctx, clientId, userId)
			if err != nil {
				render.Render(w, r, CreateStudyTemplateErrorResponse(err))
				return
			}

			render.Respond(w, r, presenter.StudyTemplateGetsResponse{
				Status: "ok",
				Data:   mapModelsResponse(studyTemplates),
			})
			return
		}

		studyTemplate, err := h.studyTemplatesUC.GetByOwnerName(r.Context(), clientId, userId, templateName)
		if err != nil {
			render.Render(w, r, CreateStudyTemplateErrorResponse(err))
			return
		}

		render.Respond(w, r, presenter.StudyTemplateGetResponse{
			Status: "ok",
			Data:   mapModelResponse(studyTemplate),
		})

	}
}

func (h *studyTemplateHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		q := r.URL.Query()

		clientId := q.Get("client")
		userId := q.Get("user")
		templateName := q.Get("template")

		_, err := h.studyTemplatesUC.DeleteByOwnerName(ctx, clientId, userId, templateName)
		if err != nil {
			render.Render(w, r, CreateStudyTemplateErrorResponse(err))
			return
		}

		render.Respond(w, r, presenter.StudyTemplateDeleteResponse{Status: "ok"})
	}
}

func mapModelResponse(exp *models.StudyTemplate) *presenter.StudyTemplateResponse {
	return &presenter.StudyTemplateResponse{
		Name:    exp.Name,
		Content: exp.Content,
	}
}

func mapModelsResponse(exp []*models.StudyTemplate) []*presenter.StudyTemplatesResponse {
	out := make([]*presenter.StudyTemplatesResponse, len(exp))
	for i, studyTemplate := range exp {
		out[i] = &presenter.StudyTemplatesResponse{
			Name: studyTemplate.Name,
		}
	}
	return out
}

func CreateStudyTemplateErrorResponse(err error) render.Renderer {
	parsedErr := httpErrors.ParseErrors(err)

	return &presenter.StudyTemplateErrorResponse{
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
