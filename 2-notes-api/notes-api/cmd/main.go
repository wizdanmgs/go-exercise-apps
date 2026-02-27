package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	delivery "notes-api/internal/delivery/http"
	"notes-api/internal/logger"
	"notes-api/internal/repository/memory"
	"notes-api/internal/usecase"
)

func main() {
	// ==== Chi Router ====
	// Initialize logger
	logg := logger.New()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Initialize infrastructure
	repo := memory.NewMemoryRepository()

	// Inject into usecase
	noteUsecase := usecase.NewNoteUsecase(repo)

	// Inject into delivery
	handler := delivery.NewNoteHandler(noteUsecase, logg)

	// Setup Router
	r := chi.NewRouter()

	// Built-in middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Routes
	r.Route("/notes", func(r chi.Router) {

		r.Post("/", handler.Create)
		r.Get("/", handler.GetAll)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetByID)
			r.Put("/", handler.Update)
			r.Delete("/", handler.Delete)
		})
	})

	// HTTP Server
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Start Server (goroutine)
	go func() {
		logger.Info("server started", "addr", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done() // Wait for signal
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", "error", err)
	} else {
		logger.Info("server shutdown gracefully")
	}

	// ===== Standard Server Mux ====
	// logg := logger.New()
	//
	// // Initialize infrastructure
	// repo := memory.NewMemoryRepository()
	//
	// // Inject into usecase
	// noteUsecase := usecase.NewNoteUsecase(repo)
	//
	// // Inject into delivery
	// handler := delivery.NewNoteHandler(noteUsecase, logg)
	//
	//
	// // Setup Router
	// mux := http.NewServeMux()
	//
	// mux.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodGet:
	// 		handler.GetAll(w, r)
	// 	case http.MethodPost:
	// 		handler.Create(w, r)
	// 	default:
	// 		http.NotFound(w, r)
	// 	}
	// })
	//
	// mux.HandleFunc("/notes/", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodGet:
	// 		handler.GetByID(w, r)
	// 	case http.MethodPut:
	// 		handler.Update(w, r)
	// 	case http.MethodDelete:
	// 		handler.Delete(w, r)
	// 	default:
	// 		http.NotFound(w, r)
	// 	}
	// })
	//
	// // Wrap with logging middleware
	// loggedMux := delivery.LoggingMiddleware(logg, mux)
	//
	// logg.Info("server_started", "port", 8080)
	//
	// if err := http.ListenAndServe(":8080", loggedMux); err != nil {
	// 	logg.Error("server_failed", "error", err)
	// 	log.Fatal(err)
	// }
}
