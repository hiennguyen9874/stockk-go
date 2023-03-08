package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hiennguyen9874/stockk-go/config"
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
}

// NewServer creates and configures an APIServer serving all application routes.
func NewServer(cfg *config.Config, db *gorm.DB) (*Server, error) {
	log.Println("configuring server...")
	api, err := New(db, cfg)

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
		cfg: cfg,
		db:  db,
	}, nil
}

// Start runs ListenAndServe on the http.Server with graceful shutdown.
func (srv *Server) Start() {
	log.Println("starting server...")

	go func() {
		log.Printf("Listening on %s\n", srv.server.Addr)
		if err := srv.server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit

	log.Println("Shutting down server... Reason:", sig)

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	if err := srv.server.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Println("Server gracefully stopped")
}
