package service

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"url-shortener/internal/model"
	"url-shortener/internal/store"
)

var ErrCollisionLimit = errors.New("unable to generate unique short code")

type Shortener struct {
	store *store.MemoryStore
}

func NewShortener(s *store.MemoryStore) *Shortener {
	return &Shortener{store: s}
}

func (s *Shortener) Create(original string, ttl time.Duration) (string, error) {
	const maxAttempts uint8 = 5

	for range maxAttempts {
		code := generateCode(original)

		if !s.store.Exists(code) {
			url := model.URL{
				Code:      code,
				Original:  original,
				ExpiresAt: time.Now().Add(ttl),
			}

			s.store.Save(url)
			return code, nil
		}
	}

	return "", ErrCollisionLimit

}

func (s *Shortener) Resolve(code string) (string, bool) {
	url, ok := s.store.Get(code)
	if !ok {
		return "", false
	}

	if time.Now().After(url.ExpiresAt) {
		s.store.Delete(code)
		return "", false
	}

	return url.Original, ok
}

func generateCode(input string) string {
	hash := sha256.Sum256([]byte(input + time.Now().String()))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return encoded[:8] // short code length
}
