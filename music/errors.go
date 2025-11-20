package music

import "errors"

var (
	// ErrNoCommits is returned when no commits are provided.
	ErrNoCommits = errors.New("no commits provided")

	// ErrInvalidMode is returned when an invalid generation mode is specified.
	ErrInvalidMode = errors.New("invalid generation mode")
)

