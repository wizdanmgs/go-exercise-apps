package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"file-upload-server/internal/handler"
	"file-upload-server/internal/service"
)

func NewRouter(uploadDir string) http.Handler {
	uploadService := service.NewUploadService(uploadDir)
	uploadHandler := handler.NewUploadHandler(uploadService)

	r := chi.NewRouter()

	// middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("multipart/form-data"))

	// Routes
	r.Route("/api", func(r chi.Router) {
		r.Post("/upload", uploadHandler.Upload)
	})
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))

	return r
}
