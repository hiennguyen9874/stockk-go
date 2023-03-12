package tickers

import (
	"github.com/hiennguyen9874/go-boilerplate/internal"
	"github.com/hiennguyen9874/go-boilerplate/internal/models"
)

type TickerUseCaseI interface {
	internal.UseCaseI[models.Ticker]
}
