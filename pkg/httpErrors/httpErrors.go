package httpErrors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrorBadRequest            = errors.New("bad request")
	ErrorWrongCredentials      = errors.New("wrong credentials")
	ErrorNotFound              = errors.New("not found")
	ErrorUnauthorized          = errors.New("unauthorized")
	ErrorForbidden             = errors.New("forbidden")
	ErrorPermissionDenied      = errors.New("permission denied")
	ErrorExpiredCSRFError      = errors.New("expired csrf token")
	ErrorWrongCSRFToken        = errors.New("wrong csrf token")
	ErrorCSRFNotPresented      = errors.New("csrf not presented")
	ErrorNotRequiredFields     = errors.New("no such required fields")
	ErrorBadQueryParams        = errors.New("invalid query params")
	ErrorInternalServerError   = errors.New("internal server error")
	ErrorRequestTimeoutError   = errors.New("request timeout")
	ErrorExistsEmailError      = errors.New("user with given email already exists")
	ErrorInvalidJWTToken       = errors.New("invalid jwt token")
	ErrorInvalidJWTClaims      = errors.New("invalid jwt claims")
	ErrorNotAllowedImageHeader = errors.New("not allowed image header")
	ErrorNoCookie              = errors.New("not found cookie header")
	ErrorValidation            = errors.New("validation")
)

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err        error  `json:"-"`             // low-level runtime error
	Status     int    `json:"status"`        // http response status code
	StatusText string `json:"statusText"`    // user-level status message
	Msg        string `json:"msg,omitempty"` // application-level error message, for debugging
}

// Error  Error() interface method
func (e ErrResponse) Error() string {
	return fmt.Sprintf("status: %d - statusText: %s - msg: %s - error: %v", e.Status, e.StatusText, e.Msg, e.Err)
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Status)
	return nil
}

func Err(err error, status int, statusText string) render.Renderer {
	return &ErrResponse{
		Err:        err,
		Status:     status,
		StatusText: statusText,
		Msg:        err.Error(),
	}
}

func ErrBadRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusBadRequest,
		StatusText: ErrorBadRequest.Error(),
		Msg:        err.Error(),
	}
}

func ErrNotFound(err error) render.Renderer {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusNotFound,
		StatusText: ErrorNotFound.Error(),
		Msg:        err.Error(),
	}
}

func ErrUnauthorized(err error) render.Renderer {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorUnauthorized.Error(),
		Msg:        err.Error(),
	}
}

func ErrForbidden(err error) render.Renderer {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusForbidden,
		StatusText: ErrorForbidden.Error(),
		Msg:        err.Error(),
	}
}

func ErrInternalServer(err error) render.Renderer {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusInternalServerError,
		StatusText: ErrorInternalServerError.Error(),
		Msg:        err.Error(),
	}
}

func ErrValidation(err error) render.Renderer {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnprocessableEntity,
		StatusText: ErrorValidation.Error(),
		Msg:        err.Error(),
	}
}
