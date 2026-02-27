package domain

import "errors"

// Domain-level error definitions.
// These are sentinel errors used across layers.

var (
	//ErrNotFound indicates entity does not exist.
	ErrNotFound = errors.New("not found")

	//ErrInvalidInput indicates validation failure.
	ErrInvalidInput = errors.New("invalid input")

	//ErrDb indicates database error.
	ErrDb = errors.New("db error")
)
