package httpErrors

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrorBadRequest          = errors.New("bad request")
	ErrorNotFound            = errors.New("not found")
	ErrorUnauthorized        = errors.New("unauthorized")
	ErrorForbidden           = errors.New("forbidden")
	ErrorInternalServerError = errors.New("internal server error")
	ErrorRequestTimeoutError = errors.New("request timeout")
	ErrorExistsEmailError    = errors.New("user with given email already exists")
	ErrorInvalidJWTToken     = errors.New("invalid jwt token")
	ErrorInvalidJWTClaims    = errors.New("invalid jwt claims")
	ErrorValidation          = errors.New("validation")
	ErrorWrongPassword       = errors.New("wrong password")
	ErrorTokenNotFound       = errors.New("token not found")
	ErrorInactiveUser        = errors.New("inactive user")
	ErrorNotEnoughPrivileges = errors.New("not enough privileges")
	ErrorGenToken            = errors.New("error when generate token")
)

// Rest error interface
type ErrRest interface {
	GetErr() error
	GetStatus() int
	GetStatusText() string
	GetMsg() string
	Error() string
}

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

func (e *ErrResponse) GetErr() error {
	return e.Err
}

func (e *ErrResponse) GetStatus() int {
	return e.Status
}

func (e *ErrResponse) GetStatusText() string {
	return e.StatusText
}

func (e *ErrResponse) GetMsg() string {
	return e.Msg
}

// Error Error() interface method
func (e *ErrResponse) Error() string {
	return fmt.Sprintf("status: %d - statusText: %s - msg: %s - error: %v", e.Status, e.StatusText, e.Msg, e.Err)
}

func ErrBadRequest(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusBadRequest,
		StatusText: ErrorBadRequest.Error(),
		Msg:        err.Error(),
	}
}

func ErrNotFound(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusNotFound,
		StatusText: ErrorNotFound.Error(),
		Msg:        err.Error(),
	}
}

func ErrUnauthorized(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorUnauthorized.Error(),
		Msg:        err.Error(),
	}
}

func ErrForbidden(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusForbidden,
		StatusText: ErrorForbidden.Error(),
		Msg:        err.Error(),
	}
}

func ErrInternalServer(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusInternalServerError,
		StatusText: ErrorInternalServerError.Error(),
		Msg:        err.Error(),
	}
}

func ErrValidation(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnprocessableEntity,
		StatusText: ErrorValidation.Error(),
		Msg:        err.Error(),
	}
}

func ErrRequestTimeoutError(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusRequestTimeout,
		StatusText: ErrorRequestTimeoutError.Error(),
		Msg:        err.Error(),
	}
}

func ErrInactiveUser(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusForbidden,
		StatusText: ErrorInactiveUser.Error(),
		Msg:        err.Error(),
	}
}

func ErrNotEnoughPrivileges(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusForbidden,
		StatusText: ErrorNotEnoughPrivileges.Error(),
		Msg:        err.Error(),
	}
}

func ErrInvalidJWTToken(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorInvalidJWTToken.Error(),
		Msg:        err.Error(),
	}
}

func ErrInvalidJWTClaims(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorInvalidJWTClaims.Error(),
		Msg:        err.Error(),
	}
}

func ErrWrongPassword(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorWrongPassword.Error(),
		Msg:        err.Error(),
	}
}

func ErrGenToken(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusBadRequest,
		StatusText: ErrorGenToken.Error(),
		Msg:        err.Error(),
	}
}

func ErrTokenNotFound(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorTokenNotFound.Error(),
		Msg:        err.Error(),
	}
}

// Parser of error string messages ,returns RestError
func ParseErrors(err error) ErrRest {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrNotFound(err)
	case errors.Is(err, context.DeadlineExceeded):
		return ErrRequestTimeoutError(err)
	case strings.Contains(err.Error(), "SQLSTATE"):
		return parseSqlErrors(err)
	default:
		if restErr, ok := err.(ErrRest); ok {
			return restErr
		}
		return ErrBadRequest(err)
	}
}

// Parser sql error, returns RestError
func parseSqlErrors(err error) ErrRest {
	if strings.Contains(err.Error(), "23505") {
		return &ErrResponse{
			Err:        err,
			Status:     http.StatusBadRequest,
			StatusText: ErrorExistsEmailError.Error(),
			Msg:        err.Error(),
		}
	}
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusBadRequest,
		StatusText: ErrorBadRequest.Error(),
		Msg:        err.Error(),
	}
}
