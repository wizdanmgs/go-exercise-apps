package http

import (
	"notes-api/internal/domain"
	"testing"
)

func TestMapErrorToStatus(t *testing.T) {
	if mapErrorToStatus(domain.ErrInvalidInput) != 400 {
		t.Fatal("wrong status for invalid input")
	}

	if mapErrorToStatus(domain.ErrNotFound) != 404 {
		t.Fatal("wrong status for not found")
	}
}
