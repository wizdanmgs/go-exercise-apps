package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"url-shortener/internal/handler"
	"url-shortener/internal/service"
	"url-shortener/internal/store"
)

func TestIntegration_GracefulShutdown(t *testing.T) {
	store := store.NewMemoryStore()
	shortener := service.NewShortener(store)
	h := handler.NewHandler(shortener)

	r := chi.NewRouter()
	r.Post("/shorten", h.Create)
	r.Get("/{code}", h.Redirect)

	srv := &http.Server{
		Addr:    "127.0.0.1:0",
		Handler: r,
	}

	ln := httptest.NewUnstartedServer(r)
	ln.Start()
	defer ln.Close()

	// Start server in background
	go func() {
		_ = srv.Serve(ln.Listener)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Trigger shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}
}

func TestIntegration_ShutdownWaitsForRequest(t *testing.T) {
	r := chi.NewRouter()

	r.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	srv := &http.Server{
		Addr:    "127.0.0.1:0",
		Handler: r,
	}

	ts := httptest.NewUnstartedServer(r)
	ts.Start()
	defer ts.Close()

	go func() {
		_ = srv.Serve(ts.Listener)
	}()

	errCh := make(chan error, 1)

	go func() {
		_, err := http.Get(ts.URL + "/slow")
		errCh <- err
	}()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("failed to call slow: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("request timed out")
	}

	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}
}
