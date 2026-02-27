package model

import "time"

type URL struct {
	Code      string
	Original  string
	ExpiresAt time.Time
}
