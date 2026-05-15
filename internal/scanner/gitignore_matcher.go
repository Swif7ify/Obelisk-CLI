package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// GitignoreMatcher handles .gitignore pattern matching for file scanning.
type GitignoreMatcher struct {
	patterns []string
	root     string
}

// NewGitignoreMatcher creates a new gitignore matcher from a project path.
func NewGitignoreMatcher(projectPath string) (*GitignoreMatcher, error) {
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	patterns, err := parseGitignoreFile(gitignorePath)
	if err != nil {
		// If .gitignore doesn't exist or can't be read, return empty matcher
		return &GitignoreMatcher{
			patterns: []string{},
			root:     projectPath,
		}, nil
	}

	return &GitignoreMatcher{
		patterns: patterns,
		root:     projectPath,
	}, nil
}

// ShouldIgnore checks if a path should be ignored based on .gitignore patterns.
func (m *GitignoreMatcher) ShouldIgnore(path string) bool {
	// Always ignore common directories regardless of .gitignore
	alwaysIgnore := []string{
		"node_modules",
		".git",
		"vendor",
		"dist",
		"build",
		".next",
		"__pycache__",
		".cache",
		"coverage",
		".github",
		".vscode",
		".idea",
	}

	// Get relative path from project root
	relPath, err := filepath.Rel(m.root, path)
	if err != nil {
		relPath = filepath.Base(path)
	}

	// Check always-ignore list
	parts := strings.Split(relPath, string(filepath.Separator))
	for _, part := range parts {
		for _, ignore := range alwaysIgnore {
			if part == ignore {
				return true
			}
		}
	}

	// Check .gitignore patterns
	for _, pattern := range m.patterns {
		if matchPattern(pattern, relPath) {
			return true
		}
	}

	return false
}

// parseGitignoreFile reads and parses a .gitignore file.
func parseGitignoreFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Skip negation patterns for simplicity
		if strings.HasPrefix(line, "!") {
			continue
		}
		patterns = append(patterns, line)
	}

	return patterns, scanner.Err()
}

// matchPattern checks if a path matches a gitignore pattern.
func matchPattern(pattern, path string) bool {
	// Normalize paths
	pattern = filepath.ToSlash(pattern)
	path = filepath.ToSlash(path)

	// Remove leading/trailing slashes
	pattern = strings.Trim(pattern, "/")
	path = strings.Trim(path, "/")

	// Handle directory patterns (ending with /)
	if strings.HasSuffix(pattern, "/") {
		pattern = strings.TrimSuffix(pattern, "/")
		// Match if path starts with pattern
		return strings.HasPrefix(path, pattern+"/") || path == pattern
	}

	// Handle wildcard patterns
	if strings.Contains(pattern, "*") {
		return matchWildcard(pattern, path)
	}

	// Exact match
	if pattern == path {
		return true
	}

	// Match if pattern matches any part of the path
	pathParts := strings.Split(path, "/")
	for i := range pathParts {
		subPath := strings.Join(pathParts[i:], "/")
		if pattern == subPath {
			return true
		}
	}

	// Match basename
	if filepath.Base(path) == pattern {
		return true
	}

	return false
}

// matchWildcard performs simple wildcard matching.
func matchWildcard(pattern, path string) bool {
	// Handle ** (match any number of directories)
	if strings.Contains(pattern, "**") {
		parts := strings.Split(pattern, "**")
		if len(parts) == 2 {
			prefix := strings.TrimSuffix(parts[0], "/")
			suffix := strings.TrimPrefix(parts[1], "/")
			
			if prefix != "" && !strings.HasPrefix(path, prefix) {
				return false
			}
			if suffix != "" && !strings.HasSuffix(path, suffix) {
				return false
			}
			return true
		}
	}

	// Handle * (match within a path segment)
	if strings.HasPrefix(pattern, "*.") {
		// Extension match (e.g., *.log)
		ext := strings.TrimPrefix(pattern, "*")
		return strings.HasSuffix(path, ext)
	}

	if strings.HasSuffix(pattern, "*") {
		// Prefix match (e.g., .env*)
		prefix := strings.TrimSuffix(pattern, "*")
		basename := filepath.Base(path)
		return strings.HasPrefix(basename, prefix)
	}

	// Simple contains check for other wildcards
	if strings.Contains(pattern, "*") {
		parts := strings.Split(pattern, "*")
		if len(parts) == 2 {
			basename := filepath.Base(path)
			return strings.HasPrefix(basename, parts[0]) && strings.HasSuffix(basename, parts[1])
		}
	}

	return false
}

// Made with Bob
