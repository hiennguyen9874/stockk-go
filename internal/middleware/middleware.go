package middleware

import (
	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
)

type MiddlewareManager struct {
	cfg    *config.Config
	logger logger.Logger
}

func CreateMiddlewareManager(cfg *config.Config, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{
		cfg:    cfg,
		logger: logger,
	}
}
