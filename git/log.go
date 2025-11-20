package git

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ParseLog parses Git log output from the specified repository path or URL.
// If repoPath is a URL, it will be cloned to a temporary directory first.
// Returns a slice of commits in chronological order (oldest first).
func ParseLog(repoPath string) ([]Commit, error) {
	isURL, err := isGitURL(repoPath)
	if err != nil {
		return nil, fmt.Errorf("invalid repository path: %w", err)
	}

	var actualPath string
	var tempDir string
	var needsCleanup bool

	if isURL {
		fmt.Printf("Cloning repository from URL: %s\n", repoPath)
		tempDir, err = os.MkdirTemp("", "git2midi-*")
		if err != nil {
			return nil, fmt.Errorf("failed to create temporary directory: %w", err)
		}
		needsCleanup = true

		cloneCmd := exec.Command("git", "clone", repoPath, tempDir)
		cloneCmd.Stderr = os.Stderr
		if err := cloneCmd.Run(); err != nil {
			os.RemoveAll(tempDir)
			return nil, fmt.Errorf("failed to clone repository: %w", err)
		}

		actualPath = tempDir
	} else {
		actualPath = repoPath
	}

	defer func() {
		if needsCleanup && tempDir != "" {
			os.RemoveAll(tempDir)
		}
	}()

	cmd := exec.Command("git", "-C", actualPath, "log", "--pretty=format:%H|%ct|%an|%s")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to read git log: %w", err)
	}

	var commits []Commit
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}

		var timestamp int64
		if ts, err := parseUnixTimestamp(parts[1]); err == nil {
			timestamp = ts
		} else {
			timestamp = time.Now().Unix()
		}

		commit := Commit{
			Hash:      parts[0],
			Timestamp: timestamp,
			Author:    parts[2],
			Message:   parts[3],
		}

		if err := commit.Validate(); err != nil {
			continue
		}

		commits = append(commits, commit)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	for i, j := 0, len(commits)-1; i < j; i, j = i+1, j-1 {
		commits[i], commits[j] = commits[j], commits[i]
	}

	return commits, nil
}

// parseUnixTimestamp parses a Unix timestamp string into an int64.
func parseUnixTimestamp(s string) (int64, error) {
	var sec int64
	_, err := fmt.Sscanf(s, "%d", &sec)
	if err != nil {
		return 0, err
	}
	return sec, nil
}

// isGitURL checks if the given string is a Git repository URL.
func isGitURL(path string) (bool, error) {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		gitDir := filepath.Join(path, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return false, nil
		}
		return false, nil
	}

	lowerPath := strings.ToLower(path)
	if strings.HasPrefix(lowerPath, "http://") ||
		strings.HasPrefix(lowerPath, "https://") ||
		strings.HasPrefix(lowerPath, "git://") ||
		strings.HasPrefix(lowerPath, "ssh://") ||
		strings.HasPrefix(path, "git@") {
		return true, nil
	}

	return false, nil
}
