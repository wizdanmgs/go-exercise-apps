package main

import (
	"log"
	"net/http"
	"url-shortener/internal/handler"
	"url-shortener/internal/service"
	"url-shortener/internal/store"
)

func main() {
	store := store.NewMemoryStore()
	shortener := service.NewShortener(store)
	handler := handler.NewHandler(shortener)

	http.HandleFunc("/shorten", handler.Create)
	http.HandleFunc("/", handler.Redirect)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
