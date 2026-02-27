package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"url-shortener/internal/service"
)

type Handler struct {
	shortener *service.Shortener
}

func NewHandler(s *service.Shortener) *Handler {
	return &Handler{shortener: s}
}

type createRequest struct {
	URL string `json:"url"`
	TTL int    `json:"ttl"` // seconds
}

type createResponse struct {
	Code string `json:"code"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	code := h.shortener.Create(req.URL, time.Duration(req.TTL)*time.Second)

	resp := createResponse{Code: code}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	original, ok := h.shortener.Resolve(code)
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, original, http.StatusFound)
}
