package usecase

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"sync"
	"time"

	"scraper/internal/domain"

	"golang.org/x/time/rate"
)

type ScraperUsecase struct {
	fetcher     domain.Fetcher
	workerCount int
	limiter     *rate.Limiter
	maxRetries  int
	baseDelay   time.Duration
}

func NewScraperUsecase(
	fetcher domain.Fetcher,
	workerCount int,
	requestsPerSecond int,
	burst int,
	maxRetries int,
	baseDelay time.Duration,
) *ScraperUsecase {
	limiter := rate.NewLimiter(
		rate.Limit(requestsPerSecond),
		burst,
	)
	return &ScraperUsecase{
		fetcher:     fetcher,
		workerCount: workerCount,
		limiter:     limiter,
		maxRetries:  maxRetries,
		baseDelay:   baseDelay,
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

			title, err := s.fetchWithRetry(ctx, j.url)
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

func (s *ScraperUsecase) fetchWithRetry(
	ctx context.Context,
	url string,
) (string, error) {
	var lastErr error

	for attempt := 0; attempt < s.maxRetries; attempt++ {
		// Rate limit per attempt
		if err := s.limiter.Wait(ctx); err != nil {
			return "", err
		}

		title, status, err := s.fetcher.FetchTitle(ctx, url)
		if err == nil {
			return title, nil
		}

		lastErr = err

		// Stop if not retryable
		if !isRetryable(status, err) {
			return "", err
		}
		// If last attempt, break
		if attempt == s.maxRetries {
			break
		}

		// Exponential backoff
		backoff := s.baseDelay * time.Duration(1<<attempt)

		// Add jitter (up to 50%)
		jitter := time.Duration(rand.Int63n(int64(backoff / 2)))
		delay := backoff + jitter

		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(delay):
		}
	}
	return "", lastErr
}

func isRetryable(status int, err error) bool {
	// Network-level errors
	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return true
		}
	}

	// Retry on context deadline exceeded? NO
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	// HTTP-level retry logic
	switch {
	case status == 429:
		return true
	case status >= 500 && status <= 599:
		return true
	}

	return false
}
