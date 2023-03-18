package presenter

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
)

type StudyTemplateCreateUpdateResponse struct {
	Status string `json:"status"`
}

type StudyTemplateGetResponse struct {
	Status string                 `json:"status"`
	Data   *StudyTemplateResponse `json:"data"`
}

type StudyTemplateResponse struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type StudyTemplatesResponse struct {
	Name string `json:"name"`
}

type StudyTemplateGetsResponse struct {
	Status string                    `json:"status"`
	Data   []*StudyTemplatesResponse `json:"data"`
}

type StudyTemplateDeleteResponse struct {
	Status string `json:"status"`
}

type StudyTemplateErrorResponse struct {
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Error   *httpErrors.ErrResponse `json:"-"`
}

func (e *StudyTemplateErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Error.Status)
	return nil
}
