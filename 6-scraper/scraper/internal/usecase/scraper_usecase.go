package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"sync"
	"time"

	"scraper/internal/domain"

	"github.com/temoto/robotstxt"
	"golang.org/x/time/rate"
)

type ScraperUsecase struct {
	fetcher       domain.Fetcher
	robotsFetcher domain.RobotsFetcher
	workerCount   int

	globalLimiter *rate.Limiter

	domainLimiters map[string]*rate.Limiter
	mu             sync.Mutex

	maxRetries int
	baseDelay  time.Duration

	domainRPS   int
	domainBurst int

	breakers         map[string]*CircuitBreaker
	breakerThreshold int
	breakerTImeout   time.Duration

	robotsCache map[string]*robotstxt.RobotsData
	robotsMu    sync.Mutex
	userAgent   string
}

func NewScraperUsecase(
	fetcher domain.Fetcher,
	robotsFetcher domain.RobotsFetcher,
	workerCount int,

	globalRPS int,
	globalBurst int,

	domainRPS int,
	domainBurst int,

	maxRetries int,
	baseDelay time.Duration,
) *ScraperUsecase {

	return &ScraperUsecase{
		fetcher:       fetcher,
		robotsFetcher: robotsFetcher,
		workerCount:   workerCount,

		globalLimiter:  rate.NewLimiter(rate.Limit(globalRPS), globalBurst),
		domainLimiters: make(map[string]*rate.Limiter),

		domainRPS:   domainRPS,
		domainBurst: domainBurst,

		maxRetries: maxRetries,
		baseDelay:  baseDelay,

		breakers:         make(map[string]*CircuitBreaker),
		breakerThreshold: 5,
		breakerTImeout:   30 * time.Second,

		robotsCache: make(map[string]*robotstxt.RobotsData),
		userAgent:   "MyCrawlerBot",
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
			if err := s.globalLimiter.Wait(ctx); err != nil {
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
	rawUrl string,
) (string, error) {
	host, err := extractHost(rawUrl)
	if err != nil {
		return "", err
	}

	group, err := s.getRobots(ctx, host)
	if err == nil && group != nil {
		allowed := group.Test(rawUrl)
		if !allowed {
			return "", fmt.Errorf("blocked by robots.txt: %s", rawUrl)
		}
	}

	if group != nil && group.CrawlDelay > 0 {
		delay := time.Duration(group.CrawlDelay) * time.Second

		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(delay):
		}
	}

	domainLimiter := s.getDomainLimiter(host)
	breaker := s.getBreaker(host)

	if !breaker.Allow() {
		return "", fmt.Errorf("circuit open for domain %s", host)
	}

	var lastErr error

	for attempt := 0; attempt < s.maxRetries; attempt++ {
		// Rate limit per attempt for Global
		if err := s.globalLimiter.Wait(ctx); err != nil {
			return "", err
		}

		// Domain-specific rate limit per attempt
		if err := domainLimiter.Wait(ctx); err != nil {
			return "", err
		}

		title, status, err := s.fetcher.FetchTitle(ctx, rawUrl)
		if err == nil {
			return title, nil
		}

		lastErr = err

		// Stop if not retryable
		if !isRetryable(status, err) {
			breaker.Failure()
			return "", err
		}
		// If last attempt, break
		if attempt == s.maxRetries {
			breaker.Failure()
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

func (s *ScraperUsecase) getDomainLimiter(host string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	limiter, exists := s.domainLimiters[host]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(s.domainRPS), s.domainBurst)
		s.domainLimiters[host] = limiter
	}

	return limiter
}

func extractHost(rawUrl string) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}

func (s *ScraperUsecase) getBreaker(host string) *CircuitBreaker {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, exists := s.breakers[host]
	if !exists {
		b = NewCircuitBreaker(s.breakerThreshold, s.breakerTImeout)
		s.breakers[host] = b
	}

	return b
}

func (s *ScraperUsecase) getRobots(
	ctx context.Context,
	host string,
) (*robotstxt.Group, error) {
	s.robotsMu.Lock()
	data, exists := s.robotsCache[host]
	s.robotsMu.Unlock()

	if !exists {
		robotsUrl := fmt.Sprintf("https://%s/robots.txt", host)

		body, err := s.robotsFetcher.FetchRobots(ctx, robotsUrl)
		if err != nil {
			return nil, err
		}

		data, err = robotstxt.FromBytes(body)
		if err != nil {
			return nil, err
		}

		s.robotsMu.Lock()
		s.robotsCache[host] = data
		s.robotsMu.Unlock()
	}

	group := data.FindGroup(s.userAgent)
	return group, nil
}
