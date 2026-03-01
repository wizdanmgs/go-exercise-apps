package repository

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

type HTTPFetcher struct {
	client *http.Client
}

type HTTPRobotsFetcher struct {
	client *http.Client
}

func NewHTTPFetcher(client *http.Client) *HTTPFetcher {
	return &HTTPFetcher{client: client}
}

func NewHTTPRobotsFetcher(client *http.Client) *HTTPRobotsFetcher {
	return &HTTPRobotsFetcher{client: client}
}

func (h *HTTPFetcher) FetchTitle(ctx context.Context, url string) (string, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", 0, err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); err == nil {
			err = closeErr
		}
	}()

	if resp.StatusCode >= 400 {
		return "", resp.StatusCode, fmt.Errorf("http error: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", resp.StatusCode, err
	}

	var title string
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	if title == "" {
		return "", resp.StatusCode, fmt.Errorf("title not found")
	}

	return title, resp.StatusCode, nil
}

func (h *HTTPRobotsFetcher) FetchRobots(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); err == nil {
			err = closeErr
		}
	}()

	return io.ReadAll(resp.Body)
}
