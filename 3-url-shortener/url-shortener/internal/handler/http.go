package handler

import (
	"encoding/json"
	"net/http"
	"time"
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
	json.NewEncoder(w).Encode(createResponse{Code: code})
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:] // /abc123 -> abc123
	original, ok := h.shortener.Resolve(code)
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, original, http.StatusFound)
}
