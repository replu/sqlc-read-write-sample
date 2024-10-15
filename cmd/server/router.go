package main

import (
	"net/http"

	"github.com/replu/sqlc-read-write-sample/internal/handler"
)

func routing(h *handler.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /user/{id}", h.Get)
	mux.HandleFunc("GET /user/tx/{id}", h.GetWithTx)
	mux.HandleFunc("POST /user", h.Create)

	return mux
}
