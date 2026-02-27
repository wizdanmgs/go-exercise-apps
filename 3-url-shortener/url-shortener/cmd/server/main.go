package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"url-shortener/internal/handler"
	"url-shortener/internal/service"
	"url-shortener/internal/store"
)

func main() {
	store := store.NewMemoryStore()
	shortener := service.NewShortener(store)
	handler := handler.NewHandler(shortener)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.AllowContentType("application/json"))

	r.Post("/shorten", handler.Create)
	r.Get("/{code}", handler.Redirect)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

