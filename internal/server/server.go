package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"task_scheduler/internal/api"
	"task_scheduler/internal/config"
)

const webDir = "./web"

type Server struct {
	logger     *log.Logger
	httpServer *http.Server
}

func New(cfg *config.ServerConfig, logger *log.Logger) *Server {
	mux := http.NewServeMux()

	mux.Handle("GET /", http.FileServer(http.Dir(webDir)))
	api.Init(mux, cfg.Password)

	return &Server{
		logger: logger,
		httpServer: &http.Server{
			Addr:        cfg.Port,
			Handler:     mux,
			ErrorLog:    logger,
			ReadTimeout: 5 * time.Second,
		},
	}
}

func (server *Server) Run() error {
	server.logger.Printf("server started: addr %s", server.httpServer.Addr)

	err := server.httpServer.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (server *Server) Shutdown(ctx context.Context) error {
	server.logger.Printf("server stopping")
	return server.httpServer.Shutdown(ctx)
}
