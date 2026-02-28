package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"scraper/internal/repository"
	"scraper/internal/usecase"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	fetcher := repository.NewHTTPFetcher(client)

	scraper := usecase.NewScraperUsecase(
		fetcher,
		5,                    // workers
		2,                    // 2 req/sec
		2,                    // burst size
		3,                    // max retries
		500*time.Millisecond, // base delay
	)

	urls := []string{
		"https://golang.org",
		"https://github.com",
		"https://google.com",
	}

	pages, _ := scraper.Crawl(ctx, urls)

	for _, p := range pages {
		fmt.Printf("%s -> %s\n", p.URL, p.Title)
	}
}
