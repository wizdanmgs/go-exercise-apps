package store

import model "url-shortener/internal/model"

type URLStore interface {
	Save(url model.URL)
	Get(code string) (model.URL, bool)
	Delete(code string)
	Exists(code string) bool
}
