package usecase

import (
	"context"
	"sync"

	"scraper/internal/domain"

	"golang.org/x/time/rate"
)

type ScraperUsecase struct {
	fetcher     domain.Fetcher
	workerCount int
	limiter     *rate.Limiter
}

func NewScraperUsecase(
	fetcher domain.Fetcher,
	workerCount int,
	requestsPerSecond int,
	burst int,
) *ScraperUsecase {
	limiter := rate.NewLimiter(
		rate.Limit(requestsPerSecond),
		burst,
	)
	return &ScraperUsecase{
		fetcher:     fetcher,
		workerCount: workerCount,
		limiter:     limiter,
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

			// Rate limiting (context-aware)
			if err := s.limiter.Wait(ctx); err != nil {
				results <- result{err: err}
				return
			}

			title, err := s.fetcher.FetchTitle(ctx, j.url)
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
