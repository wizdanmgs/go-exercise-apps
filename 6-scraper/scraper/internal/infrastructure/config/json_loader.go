package config

import (
	"encoding/json"
	"os"

	"scraper/internal/domain"
)

type JSONLoader struct {
	path string
}

func NewJSONLoader(path string) *JSONLoader {
	return &JSONLoader{path: path}
}

func (l *JSONLoader) Load() (*domain.CrawlConfig, error) {
	file, err := os.Open(l.path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var cfg domain.CrawlConfig
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
