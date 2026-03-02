package usecase

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type realFetcher struct {
	client *http.Client
}

func (r *realFetcher) FetchTitle(
	ctx context.Context,
	rawURL string,
) (string, int, error) {

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	resp, err := r.client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", resp.StatusCode,
			fmt.Errorf("http status %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	title := extractTitle(body)

	return title, resp.StatusCode, nil
}

type realRobotsFetcher struct {
	client *http.Client
}

func (r *realRobotsFetcher) FetchRobots(
	ctx context.Context,
	rawURL string,
) ([]byte, error) {

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func extractTitle(body []byte) string {
	s := string(body)

	startTag := "<title>"
	endTag := "</title>"

	start := strings.Index(s, startTag)
	if start == -1 {
		return ""
	}

	start += len(startTag)

	end := strings.Index(s[start:], endTag)
	if end == -1 {
		return ""
	}

	return s[start : start+end]
}

func TestIntegration_CrawlSuccess(t *testing.T) {
	handler := http.NewServeMux()

	handler.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User-agent: *\nAllow: /"))
	})

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<title>Hello</title>"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	fetcher := &realFetcher{client: server.Client()}
	robots := &realRobotsFetcher{client: server.Client()}

	s := NewScraperUsecase(
		fetcher,
		robots,
		2,
		100, 100,
		100, 100,
		2,
		10*time.Millisecond,
	)

	ctx := context.Background()

	pages, err := s.Crawl(ctx, []string{server.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(pages))
	}

	if pages[0].Title != "Hello" {
		t.Fatalf("unexpected title: %s", pages[0].Title)
	}
}

func TestIntegration_RobotsBlocked(t *testing.T) {
	handler := http.NewServeMux()

	handler.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User-agent: *\nDisallow: /"))
	})

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<title>Blocked</title>"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	fetcher := &realFetcher{client: server.Client()}
	robots := &realRobotsFetcher{client: server.Client()}

	s := NewScraperUsecase(
		fetcher,
		robots,
		1,
		100, 100,
		100, 100,
		1,
		10*time.Millisecond,
	)

	ctx := context.Background()

	pages, _ := s.Crawl(ctx, []string{server.URL})

	if len(pages) != 0 {
		t.Fatal("expected no pages due to robots block")
	}
}

func TestIntegration_RetryOn500(t *testing.T) {
	attempt := 0

	handler := http.NewServeMux()

	handler.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User-agent: *\nAllow: /"))
	})

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		attempt++
		if attempt < 2 {
			w.WriteHeader(500)
			return
		}
		w.Write([]byte("<title>Recovered</title>"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	fetcher := &realFetcher{client: server.Client()}
	robots := &realRobotsFetcher{client: server.Client()}

	s := NewScraperUsecase(
		fetcher,
		robots,
		1,
		100, 100,
		100, 100,
		3,
		10*time.Millisecond,
	)

	ctx := context.Background()

	pages, err := s.Crawl(ctx, []string{server.URL})
	if err != nil {
		t.Fatal(err)
	}

	if pages[0].Title != "Recovered" {
		t.Fatal("retry did not recover")
	}
}

func TestIntegration_CircuitBreaker(t *testing.T) {
	handler := http.NewServeMux()

	handler.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User-agent: *\nAllow: /"))
	})

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	fetcher := &realFetcher{client: server.Client()}
	robots := &realRobotsFetcher{client: server.Client()}

	s := NewScraperUsecase(
		fetcher,
		robots,
		1,
		100, 100,
		100, 100,
		5,
		10*time.Millisecond,
	)

	ctx := context.Background()

	// Repeated calls should eventually open circuit
	for i := 0; i < 10; i++ {
		s.fetchWithRetry(ctx, server.URL)
	}

	_, err := s.fetchWithRetry(ctx, server.URL)
	if err == nil {
		t.Fatal("expected circuit breaker open error")
	}
}
