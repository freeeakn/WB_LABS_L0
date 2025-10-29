package handler

import (
	"backend/internal/cache"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func New(c *cache.Cache) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/order/{uid}", getOrder(c)).Methods(http.MethodGet)
	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)
	return r
}

func getOrder(c *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		uid := mux.Vars(req)["uid"]
		o, ok := c.Get(uid)
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(o)
	}
}
