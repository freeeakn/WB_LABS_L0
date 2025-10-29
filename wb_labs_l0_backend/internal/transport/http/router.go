package httptransport

import (
	"context"
	"net/http"
	"time"

	"wb_labs_l0_backend/internal/config"
	"wb_labs_l0_backend/internal/logger"
	"wb_labs_l0_backend/internal/service"

	"github.com/gorilla/mux"
)

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
	svc        *service.OrderService
	log        logger.Logger
}

func NewServer(cfg *config.Config, svc *service.OrderService, log logger.Logger) *Server {
	r := mux.NewRouter()
	h := NewHandler(svc, log)

	r.HandleFunc("/orders/{id}", h.GetOrderHandler).Methods("GET")
	r.HandleFunc("/health", h.Health).Methods("GET")
	// r.PathPrefix("/").Handler(http.FileServer(http.Dir("frontend/dist/")))

	srv := &http.Server{
		Addr:         cfg.HTTPPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return &Server{
		httpServer: srv,
		cfg:        cfg,
		svc:        svc,
		log:        log,
	}
}

func (s *Server) Start() error {
	s.log.Infof("http server listening on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) {
	_ = s.httpServer.Shutdown(ctx)
}
