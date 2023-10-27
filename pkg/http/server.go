package http

import (
	"context"
	"github.com/mxmrykov/L0/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Start(cfg *config.HTTP) error {
	s.httpServer = &http.Server{
		Addr: ":" + cfg.Port,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
