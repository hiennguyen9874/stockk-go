package repository

import (
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/repository"
	"github.com/hiennguyen9874/stockk-go/internal/tickers"
	"github.com/redis/go-redis/v9"
)

type TickerRedisRepo struct {
	repository.RedisRepo[models.Ticker]
}

func CreateTickerRedisRepository(redisClient *redis.Client) tickers.TickerRedisRepository {
	return &TickerRedisRepo{
		RedisRepo: repository.RedisRepo[models.Ticker](repository.CreateRedisRepo[models.Ticker](redisClient)),
	}
}
