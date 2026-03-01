package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"scraper/internal/infrastructure/config"
	"scraper/internal/repository"
	"scraper/internal/usecase"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Load config
	loader := config.NewJSONLoader("config.json")
	cfg, err := loader.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	fetcher := repository.NewHTTPFetcher(client)
	robotsFetcher := repository.NewHTTPRobotsFetcher(client)

	scraper := usecase.NewScraperUsecase(
		fetcher,
		robotsFetcher,
		10, // workers

		5, 5, // global: 5 rps, burst 5
		2, 2, // per-domain: 2 rps, burst 2

		3,                    // max retries
		500*time.Millisecond, // base delay
	)

	pages, err := scraper.Crawl(ctx, cfg.URLs)
	if err != nil {
		log.Fatalf("crawl failed: %v", err)
	}

	for _, p := range pages {
		fmt.Printf("%s -> %s\n", p.URL, p.Title)
	}
}
