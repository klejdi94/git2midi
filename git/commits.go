package git

import (
	"errors"
	"fmt"
)

// Commit represents a parsed Git commit.
type Commit struct {
	Hash      string
	Timestamp int64
	Author    string
	Message   string
}

// Validate validates the commit data.
func (c *Commit) Validate() error {
	if c.Hash == "" {
		return errors.New("commit hash cannot be empty")
	}
	if len(c.Hash) < 7 {
		return fmt.Errorf("commit hash too short: %s", c.Hash)
	}
	return nil
}

// LimitCommits limits the number of commits to the specified count.
func LimitCommits(commits []Commit, count int) []Commit {
	if count <= 0 || count >= len(commits) {
		return commits
	}
	return commits[:count]
}

// SampleCommits evenly samples commits from the slice.
func SampleCommits(commits []Commit, count int) ([]Commit, error) {
	if count <= 0 {
		return nil, errors.New("sample count must be positive")
	}
	if count >= len(commits) {
		return commits, nil
	}

	sampled := make([]Commit, 0, count)
	step := float64(len(commits)) / float64(count)

	for i := 0; i < count; i++ {
		index := int(float64(i) * step)
		if index < len(commits) {
			sampled = append(sampled, commits[index])
		}
	}

	return sampled, nil
}

