package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/replu/sqlc-read-write-sample/internal/service"
)

type Handler struct {
	srv *service.Service
}

func NewHandler(srv *service.Service) *Handler {
	return &Handler{
		srv: srv,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "Start Get")

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.ErrorContext(r.Context(), "invalid id", slog.String("id", idStr))
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	slog.InfoContext(r.Context(), "Get user", slog.Int64("id", id))
	user, err := h.srv.Get(r.Context(), id)
	if err != nil {
		slog.ErrorContext(r.Context(), "service error", slog.String("error", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		slog.ErrorContext(r.Context(), "user not found", slog.Int64("id", id))
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// ok
	slog.InfoContext(r.Context(), "User found", slog.Any("user", user))
	if err := json.NewEncoder(w).Encode(user); err != nil {
		slog.ErrorContext(r.Context(), "json encode error", slog.String("error", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetWithTx(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "Start Get with tx")

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.ErrorContext(r.Context(), "invalid id", slog.String("id", idStr))
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	slog.InfoContext(r.Context(), "Get user with tx", slog.Int64("id", id))
	user, err := h.srv.GetWithTx(r.Context(), id)
	if err != nil {
		slog.ErrorContext(r.Context(), "service error", slog.String("error", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		slog.ErrorContext(r.Context(), "user not found", slog.Int64("id", id))
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// ok
	slog.InfoContext(r.Context(), "User found", slog.Any("user", user))
	if err := json.NewEncoder(w).Encode(user); err != nil {
		slog.ErrorContext(r.Context(), "json encode error", slog.String("error", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type CreateUser struct {
	Name string `json:"name"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "Start Create")

	var input CreateUser
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.ErrorContext(r.Context(), "invalid request", slog.String("error", err.Error()))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if input.Name == "" {
		slog.ErrorContext(r.Context(), "empty name")
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	slog.InfoContext(r.Context(), "Create user", slog.Any("user", input))

	user, err := h.srv.Create(r.Context(), input.Name)
	if err != nil {
		slog.ErrorContext(r.Context(), "service error", slog.String("error", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// ok
	slog.InfoContext(r.Context(), "User create", slog.Any("user", user))
	if err := json.NewEncoder(w).Encode(user); err != nil {
		slog.ErrorContext(r.Context(), "json encode error", slog.String("error", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
