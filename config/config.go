package config

import (
	"errors"
	"fmt"
)

// Config holds all configuration for the MIDI generation process.
type Config struct {
	RepoPath   string
	OutputPath string
	BPM        int
	Ticks      int
	Duration   int
	MaxCommits int
	Sample     bool
	Mode       Mode
}

// Mode represents the generation mode.
type Mode int

const (
	// ModeSingleTrack generates a single MIDI track with all commits.
	ModeSingleTrack Mode = iota

	// ModePerAuthor generates separate tracks for each author.
	ModePerAuthor
)

// String returns the string representation of the mode.
func (m Mode) String() string {
	switch m {
	case ModeSingleTrack:
		return "single-track"
	case ModePerAuthor:
		return "per-author"
	default:
		return "unknown"
	}
}

// ParseMode parses a mode string into a Mode value.
func ParseMode(s string) (Mode, error) {
	switch s {
	case "single-track":
		return ModeSingleTrack, nil
	case "per-author":
		return ModePerAuthor, nil
	default:
		return ModeSingleTrack, fmt.Errorf("invalid mode: %s (must be 'single-track' or 'per-author')", s)
	}
}

const (
	// DefaultBPM is the default tempo in beats per minute.
	DefaultBPM = 140

	// DefaultTicks is the default ticks per quarter note.
	DefaultTicks = 480

	// DefaultDuration is the default note duration in ticks.
	DefaultDuration = 120

	// DefaultOutputPath is the default output file path.
	DefaultOutputPath = "commits.mid"

	// DefaultRepoPath is the default repository path.
	DefaultRepoPath = "."

	// MinBPM is the minimum allowed BPM.
	MinBPM = 20

	// MaxBPM is the maximum allowed BPM.
	MaxBPM = 300

	// MinTicks is the minimum allowed ticks per quarter note.
	MinTicks = 96

	// MaxTicks is the maximum allowed ticks per quarter note.
	MaxTicks = 960

	// MinDuration is the minimum allowed note duration in ticks.
	MinDuration = 1

	// MaxDuration is the maximum allowed note duration in ticks.
	MaxDuration = 4800

	// RecommendedMaxCommits is the recommended maximum commits for reasonable file size.
	RecommendedMaxCommits = 2000
)

// Validate validates the configuration and returns an error if invalid.
func (c *Config) Validate() error {
	if c.RepoPath == "" {
		return errors.New("repository path cannot be empty")
	}

	if c.OutputPath == "" {
		return errors.New("output path cannot be empty")
	}

	if c.BPM < MinBPM || c.BPM > MaxBPM {
		return fmt.Errorf("BPM must be between %d and %d, got %d", MinBPM, MaxBPM, c.BPM)
	}

	if c.Ticks < MinTicks || c.Ticks > MaxTicks {
		return fmt.Errorf("ticks must be between %d and %d, got %d", MinTicks, MaxTicks, c.Ticks)
	}

	if c.Duration < MinDuration || c.Duration > MaxDuration {
		return fmt.Errorf("duration must be between %d and %d, got %d", MinDuration, MaxDuration, c.Duration)
	}

	if c.MaxCommits < 0 {
		return fmt.Errorf("max commits cannot be negative, got %d", c.MaxCommits)
	}

	return nil
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	return &Config{
		RepoPath:   DefaultRepoPath,
		OutputPath: DefaultOutputPath,
		BPM:        DefaultBPM,
		Ticks:      DefaultTicks,
		Duration:   DefaultDuration,
		MaxCommits: 0,
		Sample:     false,
		Mode:       ModeSingleTrack,
	}
}
