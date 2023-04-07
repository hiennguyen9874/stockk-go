package repository

import (
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/repository"
	"github.com/hiennguyen9874/stockk-go/internal/stockssnapshot"
	"github.com/redis/go-redis/v9"
)

type StockSnapshotRedisRepo struct {
	repository.RedisRepo[models.StockSnapshot]
}

func CreateStockSnapshotRedisRepository(redisClient *redis.Client) stockssnapshot.StockSnapshotRedisRepository {
	return &StockSnapshotRedisRepo{
		RedisRepo: repository.RedisRepo[models.StockSnapshot](repository.CreateRedisRepo[models.StockSnapshot](redisClient)),
	}
}
