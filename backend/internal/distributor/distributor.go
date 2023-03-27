package distributor

import (
	"github.com/hibiken/asynq"
	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
)

type RedisTaskDistributor struct {
	RedisClient *asynq.Client
	Cfg         *config.Config
	Logger      logger.Logger
}

func NewRedisClient(cfg *config.Config) *asynq.Client {
	redisOpt := asynq.RedisClientOpt{
		Addr: cfg.TaskRedis.TaskRedisAddr,
		DB:   cfg.TaskRedis.TaskRedisDb,
	}

	return asynq.NewClient(redisOpt)
}

func NewRedisTaskDistributor(redisClient *asynq.Client, cfg *config.Config, loggger logger.Logger) RedisTaskDistributor {
	return RedisTaskDistributor{
		RedisClient: redisClient,
		Cfg:         cfg,
		Logger:      loggger,
	}
}
