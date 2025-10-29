package httptransport

import (
	"context"
	"encoding/json"
	"net/http"

	"wb_labs_l0_backend/internal/logger"
	"wb_labs_l0_backend/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	svc *service.OrderService
	log logger.Logger
}

func NewHandler(svc *service.OrderService, log logger.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (h *Handler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	o, err := h.svc.GetOrder(ctx, id)
	if err != nil {
		h.log.Warnf("get order failed: %v", err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(o)
}
