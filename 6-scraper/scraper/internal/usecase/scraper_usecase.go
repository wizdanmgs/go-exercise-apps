package usecase

import (
	"context"
	"sync"
	"time"

	"scraper/internal/domain"
)

type ScraperUsecase struct {
	fetcher     domain.Fetcher
	workerCount int
	rateLimiter <-chan time.Time
}

func NewScraperUsecase(
	fetcher domain.Fetcher,
	workerCount int,
	requestsPerSecond int,
) *ScraperUsecase {
	return &ScraperUsecase{
		fetcher:     fetcher,
		workerCount: workerCount,
		rateLimiter: time.Tick(time.Second / time.Duration(requestsPerSecond)),
	}
}

type job struct {
	url string
}

type result struct {
	page domain.Page
	err  error
}

func (s *ScraperUsecase) Crawl(ctx context.Context, urls []string) ([]domain.Page, error) {
	jobs := make(chan job)
	results := make(chan result)

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < s.workerCount; i++ {
		wg.Add(1)
		go s.worker(ctx, &wg, jobs, results)
	}

	// Send jobs
	go func() {
		for _, u := range urls {
			jobs <- job{url: u}
		}
		close(jobs)
	}()

	// Close results when done
	go func() {
		wg.Wait()
		close(results)
	}()

	var pages []domain.Page

	for r := range results {
		if r.err != nil {
			continue
		}
		pages = append(pages, r.page)
	}

	return pages, nil
}

func (s *ScraperUsecase) worker(
	ctx context.Context,
	wg *sync.WaitGroup,
	jobs <-chan job,
	results chan<- result,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case j, ok := <-jobs:
			if !ok {
				return
			}

			<-s.rateLimiter

			title, err := s.fetcher.FetchTitle(j.url)
			if err != nil {
				results <- result{err: err}
				continue
			}

			results <- result{
				page: domain.Page{
					URL:   j.url,
					Title: title,
				},
			}
		}
	}
}
