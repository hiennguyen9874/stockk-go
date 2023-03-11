package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

// Server provides an http.Server.
type Server struct {
	server *http.Server
	cfg    *config.Config
	db     *gorm.DB
	logger logger.Logger
}

// NewServer creates and configures an APIServer serving all application routes.
func NewServer(cfg *config.Config, db *gorm.DB, redisClient *redis.Client, logger logger.Logger) (*Server, error) {
	logger.Info("configuring server...")
	api, err := New(db, redisClient, cfg, logger)

	if err != nil {
		return nil, err
	}

	return &Server{
		server: &http.Server{
			Addr:           cfg.Server.Port,
			Handler:        api,
			ReadTimeout:    time.Second * cfg.Server.ReadTimeout,
			WriteTimeout:   time.Second * cfg.Server.WriteTimeout,
			MaxHeaderBytes: maxHeaderBytes,
		},
		cfg:    cfg,
		db:     db,
		logger: logger,
	}, nil
}

// Start runs ListenAndServe on the http.Server with graceful shutdown.
func (srv *Server) Start() {
	srv.logger.Info("starting server...")

	go func() {
		srv.logger.Infof("Listening on %s\n", srv.server.Addr)
		if err := srv.server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit

	srv.logger.Infof("Shutting down server... Reason: %s", sig)

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	if err := srv.server.Shutdown(ctx); err != nil {
		panic(err)
	}
	srv.logger.Info("Server gracefully stopped")
}
