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
		10, // workers

		5, 5, // global: 5 rps, burst 5
		2, 2, // per-domain: 2 rps, burst 2

		3,                    // max retries
		500*time.Millisecond, // base delay
	)
	urls := []string{
		"https://golang.org",
		"https://golang.org",
		"https://wikiless.tiekoetter.com/",
		"https://wikiless.tiekoetter.com/",
		"https://wikiless.tiekoetter.com/",
		"https://example.com",
		"https://example.com",
		"https://example.com",
	}

	pages, _ := scraper.Crawl(ctx, urls)

	for _, p := range pages {
		fmt.Printf("%s -> %s\n", p.URL, p.Title)
	}
}
