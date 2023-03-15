package repository

import (
	"github.com/hiennguyen9874/stockk-go/internal/bars"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/repository"
	"github.com/redis/go-redis/v9"
)

type BarRedisRepo struct {
	repository.RedisRepo[models.Bar]
}

func CreateBarRedisRepository(redisClient *redis.Client) bars.BarRedisRepository {
	return &BarRedisRepo{
		RedisRepo: repository.RedisRepo[models.Bar](repository.CreateRedisRepo[models.Bar](redisClient)),
	}
}
