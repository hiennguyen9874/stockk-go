package presenter

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
)

type ChartCreate struct {
	OwnerSource string `json:"client" validate:"required"`
	OwnerId     string `json:"user" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Symbol      string `json:"symbol" validate:"required"`
	Resolution  string `json:"resolution" validate:"required"`
	Content     string `json:"content" validate:"required"`
}

type ChartUpdate struct {
	OwnerSource string `json:"client" validate:"required"`
	OwnerId     string `json:"user" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Symbol      string `json:"symbol" validate:"required"`
	Resolution  string `json:"resolution" validate:"required"`
	Content     string `json:"content" validate:"required"`
}

type ChartsResponse struct {
	Id         uint   `json:"id"`
	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
	Resolution string `json:"resolution"`
	Timestamp  int64  `json:"timestamp"`
}

type ChartResponse struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
	Content   string `json:"content"`
}

type ChartCreateResponse struct {
	Status string `json:"status"`
	Id     uint   `json:"id"`
}

type ChartUpdateResponse struct {
	Status string `json:"status"`
}

type ChartGetResponse struct {
	Status string         `json:"status"`
	Data   *ChartResponse `json:"data"`
}

type ChartGetsResponse struct {
	Status string            `json:"status"`
	Data   []*ChartsResponse `json:"data"`
}

type ChartDeleteResponse struct {
	Status string `json:"status"`
}

type ChartErrorResponse struct {
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Error   *httpErrors.ErrResponse `json:"-"`
}

func (e *ChartErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Error.Status)
	return nil
}
