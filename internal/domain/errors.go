package domain

import "errors"

var (
	ErrNotFound     = errors.New("data not found")
	ErrConflict     = errors.New("data already exists")
	ErrUnauthorized = errors.New("invalid credentials")
	ErrForbidden    = errors.New("access forbidden")
	ErrInvalidInput = errors.New("invalid input")
)
