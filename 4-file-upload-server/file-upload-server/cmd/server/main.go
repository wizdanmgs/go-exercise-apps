package main

import (
	"log"
	"net/http"
	"os"

	"file-upload-server/internal/server"
)

func main() {
	uploadDir := "uploads"

	// ensure uploads directory exists
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	r := server.NewRouter(uploadDir)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
