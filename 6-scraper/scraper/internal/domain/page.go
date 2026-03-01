package domain

import "context"

type Page struct {
	URL   string
	Title string
}

type Fetcher interface {
	FetchTitle(ctx context.Context, url string) (string, int, error)
}

type RobotsFetcher interface {
	FetchRobots(ctx context.Context, robotsURL string) ([]byte, error)
}
