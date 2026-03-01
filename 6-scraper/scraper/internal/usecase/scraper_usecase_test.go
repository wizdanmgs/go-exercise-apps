package usecase

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockFetcher struct {
	fetchFunc func(ctx context.Context, url string) (string, int, error)
}

func (m *mockFetcher) FetchTitle(
	ctx context.Context,
	url string,
) (string, int, error) {
	return m.fetchFunc(ctx, url)
}

type mockRobotsFetcher struct {
	body []byte
	err  error
}

func (m *mockRobotsFetcher) FetchRobots(
	ctx context.Context,
	url string,
) ([]byte, error) {
	return m.body, m.err
}

func TestFetchWithRetry_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		err         error
		maxRetries  int
		expectError bool
	}{
		{
			name:        "success",
			status:      200,
			err:         nil,
			maxRetries:  3,
			expectError: false,
		},
		{
			name:        "retryable 500",
			status:      500,
			err:         errors.New("server error"),
			maxRetries:  2,
			expectError: true,
		},
		{
			name:        "retryable 429",
			status:      429,
			err:         errors.New("too many requests"),
			maxRetries:  2,
			expectError: true,
		},
		{
			name:        "non retryable 400",
			status:      400,
			err:         errors.New("bad request"),
			maxRetries:  3,
			expectError: true,
		},
		{
			name:        "context canceled",
			status:      0,
			err:         context.Canceled,
			maxRetries:  3,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount := 0

			fetcher := &mockFetcher{
				fetchFunc: func(ctx context.Context, url string) (string, int, error) {
					callCount++
					if tt.err == nil {
						return "OK", 200, nil
					}
					return "", tt.status, tt.err
				},
			}

			robots := &mockRobotsFetcher{
				body: []byte("User-agent: *\nAllow: /"),
			}

			s := NewScraperUsecase(
				fetcher,
				robots,
				1,        // workers
				100, 100, // global limiter (very high to avoid blocking)
				100, 100, // domain limiter
				tt.maxRetries,
				1*time.Millisecond, // small delay for fast tests
			)

			ctx := context.Background()

			_, err := s.fetchWithRetry(ctx, "https://example.com")

			if tt.expectError && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.err != nil && callCount == 0 {
				t.Fatalf("expected retries but fetch not called")
			}
		})
	}
}

func TestCrawl_Success(t *testing.T) {
	fetcher := &mockFetcher{
		fetchFunc: func(ctx context.Context, url string) (string, int, error) {
			return "Title-" + url, 200, nil
		},
	}

	robots := &mockRobotsFetcher{
		body: []byte("User-agent: *\nAllow: /"),
	}

	s := NewScraperUsecase(
		fetcher,
		robots,
		3,
		100, 100,
		100, 100,
		1,
		1*time.Millisecond,
	)

	urls := []string{
		"https://a.com",
		"https://b.com",
		"https://c.com",
	}

	ctx := context.Background()

	pages, err := s.Crawl(ctx, urls)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pages) != 3 {
		t.Fatalf("expected 3 pages, got %d", len(pages))
	}
}

func TestFetchWithRetry_RobotsBlocked(t *testing.T) {
	fetcher := &mockFetcher{
		fetchFunc: func(ctx context.Context, url string) (string, int, error) {
			return "SHOULD NOT HAPPEN", 200, nil
		},
	}

	robots := &mockRobotsFetcher{
		body: []byte("User-agent: *\nDisallow: /"),
	}

	s := NewScraperUsecase(
		fetcher,
		robots,
		1,
		100, 100,
		100, 100,
		1,
		1*time.Millisecond,
	)

	ctx := context.Background()

	_, err := s.fetchWithRetry(ctx, "https://blocked.com")

	if err == nil {
		t.Fatal("expected robots.txt block error")
	}
}
