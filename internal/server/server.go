package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/pustato/image-previewer/internal/app"
	"github.com/pustato/image-previewer/internal/logger"
)

type Server struct {
	server *http.Server
}

func New(addr string, app app.App, logg logger.Logger) *Server {
	return &Server{
		server: &http.Server{
			Addr:    addr,
			Handler: &Handler{app, logg},
		},
	}
}

func (s *Server) Start() error {
	if err := s.server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server start: %w", err)
		}
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	return nil
}
