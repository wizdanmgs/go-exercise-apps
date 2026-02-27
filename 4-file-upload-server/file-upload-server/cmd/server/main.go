package main

import (
	"log"
	"net/http"
	"os"

	"file-upload-server/internal/handler"
	"file-upload-server/internal/service"
)

func main() {
	// ensure uploads directory exists
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	uploadService := service.NewUploadService("uploads")
	uploadHandler := handler.NewUploadHandler(uploadService)

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", uploadHandler.Upload)
	mux.Handle("/uploads/", http.StripPrefix("/uploads", http.FileServer(http.Dir("uploads"))))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
