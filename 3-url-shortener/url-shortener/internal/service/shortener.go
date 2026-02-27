package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"url-shortener/internal/model"
	"url-shortener/internal/store"
)

var ErrCollisionLimit = errors.New("unable to generate unique short code")

type CodeGenerator func() string

type Shortener struct {
	store store.URLStore
	gen   CodeGenerator
}

type Option func(*Shortener)

func WithGenerator(gen CodeGenerator) Option {
	return func(s *Shortener) {
		s.gen = gen
	}
}

func defaultGenerator() string {
	return generateCode()
}

func NewShortener(store store.URLStore, opts ...Option) *Shortener {
	s := &Shortener{
		store: store,
		gen:   defaultGenerator,
	}

	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Shortener) Create(original string, ttl time.Duration) (string, error) {
	const maxAttempts uint8 = 5

	for range maxAttempts {
		code := s.gen()

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

func generateCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:8]
}
