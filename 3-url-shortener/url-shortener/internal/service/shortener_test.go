package service

import (
	"testing"
	"time"
	"url-shortener/internal/store"
)

func TestShortener_Create(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		ttl     time.Duration
		wantErr bool
	}{
		{
			name:    "success",
			url:     "https://google.com",
			ttl:     time.Hour,
			wantErr: false,
		},
		{
			name:    "empty url still generates code",
			url:     "",
			ttl:     time.Hour,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := store.NewMemoryStore()
			s := NewShortener(store)

			code, err := s.Create(tt.url, tt.ttl)

			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}

			if code == "" {
				t.Fatalf("expected non-empty code")
			}

			// ensure saved
			if !store.Exists(code) {
				t.Fatalf("code was not saved in store")
			}
		})
	}
}

func TestShortener_Resolve(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(s *Shortener) string
		wantFound bool
	}{
		{
			name: "found and not expired",
			setup: func(s *Shortener) string {
				code, _ := s.Create("https://example.com", time.Hour)
				return code
			},
			wantFound: true,
		},
		{
			name: "not found",
			setup: func(s *Shortener) string {
				return "doesnotexist"
			},
			wantFound: false,
		},
		{
			name: "expired",
			setup: func(s *Shortener) string {
				code, _ := s.Create("https://expired.com", -time.Hour)
				return code
			},
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := store.NewMemoryStore()
			s := NewShortener(store)

			code := tt.setup(s)

			_, found := s.Resolve(code)

			if found != tt.wantFound {
				t.Fatalf("expected found=%v, got %v", tt.wantFound, found)
			}
		})
	}
}

func TestShortener_Create_Collision(t *testing.T) {
	store := store.NewMemoryStore()

	// always return same code to force collision
	s := NewShortener(store, WithGenerator(func() string {
		return "fixed"
	}))
	// first create works
	_, err := s.Create("url1", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// second create should fail due to collision limit
	_, err = s.Create("url2", time.Hour)
	if err == nil {
		t.Fatalf("expected collision error")
	}
}
