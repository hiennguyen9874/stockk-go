package worker

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/hiennguyen9874/stockk-go/config"

	"github.com/hiennguyen9874/stockk-go/internal/users"
	userProcessor "github.com/hiennguyen9874/stockk-go/internal/users/processor"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/hiennguyen9874/stockk-go/pkg/sendEmail"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor struct {
	server *asynq.Server
	cfg    *config.Config
	logger logger.Logger
}

func NewTaskProcessor(cfg *config.Config, logger logger.Logger) (*TaskProcessor, error) {
	redisOpt := asynq.RedisClientOpt{
		Addr: cfg.TaskRedis.TaskRedisAddr,
		DB:   cfg.TaskRedis.TaskRedisDb,
	}

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				fmt.Printf("Err: %v, Type: %v, Payload: %v, Msg: process task failed", err, task.Type(), task.Payload())
			}),
			Logger: logger,
		},
	)

	return &TaskProcessor{
		server: server,
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (taskProcessor *TaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	emailSender := sendEmail.NewEmailSender(taskProcessor.cfg)

	// Processor
	userRedisTaskProcessor := userProcessor.NewUserRedisTaskProcessor(taskProcessor.server, taskProcessor.cfg, taskProcessor.logger, emailSender)

	mux.HandleFunc(users.TaskSendEmail, userRedisTaskProcessor.ProcessTaskSendEmail)

	return taskProcessor.server.Run(mux)
}
