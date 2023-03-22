package tickers

import (
	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type TickerRedisRepository interface {
	internal.RedisRepository[models.Ticker]
}
