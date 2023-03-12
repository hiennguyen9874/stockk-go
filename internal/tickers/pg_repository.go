package tickers

import (
	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type TickerPgRepository interface {
	internal.PgRepository[models.Ticker]
}
