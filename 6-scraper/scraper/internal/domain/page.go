package domain

type Page struct {
	URL   string
	Title string
}

type Fetcher interface {
	FetchTitle(url string) (string, error)
}
