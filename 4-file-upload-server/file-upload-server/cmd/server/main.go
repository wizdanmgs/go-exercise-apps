package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"file-upload-server/internal/handler"
	"file-upload-server/internal/service"
)

func main() {
	uploadDir := "uploads"

	// ensure uploads directory exists
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

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

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
