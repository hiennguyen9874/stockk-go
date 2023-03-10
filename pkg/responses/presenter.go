package responses

import "github.com/hiennguyen9874/stockk-go/pkg/httpErrors"

type Response struct {
	Data      interface{}             `json:"data,omitempty"`
	Error     *httpErrors.ErrResponse `json:"error,omitempty"`
	IsSuccess bool                    `json:"is_success"`
}
