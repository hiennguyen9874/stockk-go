package presenter

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
)

type DrawingTemplateCreateUpdateResponse struct {
	Status string `json:"status"`
}

type DrawingTemplateResponse struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type DrawingTemplateGetResponse struct {
	Status string                   `json:"status"`
	Data   *DrawingTemplateResponse `json:"data"`
}

type DrawingTemplateGetsResponse struct {
	Status string    `json:"status"`
	Data   []*string `json:"data"`
}

type DrawingTemplateDeleteResponse struct {
	Status string `json:"status"`
}

type DrawingTemplateErrorResponse struct {
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Error   *httpErrors.ErrResponse `json:"-"`
}

func (e *DrawingTemplateErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Error.Status)
	return nil
}
