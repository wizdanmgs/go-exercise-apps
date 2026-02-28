package domain

import "context"

type Page struct {
	URL   string
	Title string
}

type Fetcher interface {
	FetchTitle(ctx context.Context, url string) (string, error)
}
