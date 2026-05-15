package scanner

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// requiredIgnorePatterns are patterns that should exist in .gitignore for security.
var requiredIgnorePatterns = map[string]string{
	".env":          "Environment files contain secrets and should be git-ignored",
	".env.local":    "Local environment files should be git-ignored",
	"node_modules/": "Dependencies should not be committed",
	".DS_Store":     "OS-specific files should be git-ignored",
	"*.log":         "Log files may contain sensitive information",
}

// sensitiveFiles are files that should NEVER be committed.
var sensitiveFiles = []string{
	".env",
	".env.local",
	".env.production",
	".env.staging",
	".env.development",
	"id_rsa",
	"id_dsa",
	"id_ecdsa",
	"id_ed25519",
	".npmrc",
	".pypirc",
}

// ScanGitignore validates .gitignore configuration and checks for tracked sensitive files.
func ScanGitignore(projectPath string) ([]Finding, error) {
	var findings []Finding

	// Check if .gitignore exists
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		findings = append(findings, Finding{
			Category:    CategorySecurity,
			Severity:    SeverityError,
			Title:       "Missing .gitignore file",
			Description: "No .gitignore file found in the project root. Sensitive files may be committed to version control",
			File:        ".gitignore",
			Suggestion:  "Create a .gitignore file with appropriate patterns for your framework",
		})
		// If no .gitignore, check for sensitive files directly
		findings = append(findings, checkSensitiveFiles(projectPath)...)
		return findings, nil
	}

	// Parse .gitignore and check for required patterns
	patterns, err := parseGitignore(gitignorePath)
	if err != nil {
		return findings, err
	}

	for pattern, reason := range requiredIgnorePatterns {
		if !containsPattern(patterns, pattern) {
			findings = append(findings, Finding{
				Category:    CategorySecurity,
				Severity:    SeverityWarning,
				Title:       "Missing .gitignore pattern: " + pattern,
				Description: reason,
				File:        ".gitignore",
				Suggestion:  "Add '" + pattern + "' to your .gitignore file",
			})
		}
	}

	// Check if any sensitive files are tracked by git
	findings = append(findings, checkTrackedSensitiveFiles(projectPath)...)

	// Check for sensitive files that exist but might not be ignored
	findings = append(findings, checkSensitiveFiles(projectPath)...)

	return findings, nil
}

// parseGitignore reads and parses a .gitignore file.
func parseGitignore(path string) ([]string, error) {
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
		patterns = append(patterns, line)
	}

	return patterns, scanner.Err()
}

// containsPattern checks if a pattern (or a close variant) exists in the patterns list.
func containsPattern(patterns []string, target string) bool {
	cleanTarget := strings.TrimSuffix(strings.TrimPrefix(target, "/"), "/")
	for _, p := range patterns {
		cleanP := strings.TrimSuffix(strings.TrimPrefix(p, "/"), "/")
		// Direct match
		if cleanP == cleanTarget {
			return true
		}
		// Wildcard match (e.g., "*.env" covers ".env")
		if strings.HasPrefix(cleanP, "*") && strings.HasSuffix(cleanTarget, strings.TrimPrefix(cleanP, "*")) {
			return true
		}
		// Pattern covers target (e.g., ".env*" covers ".env.local")
		if strings.HasSuffix(cleanP, "*") && strings.HasPrefix(cleanTarget, strings.TrimSuffix(cleanP, "*")) {
			return true
		}
	}
	return false
}

// checkTrackedSensitiveFiles uses git to check if sensitive files are being tracked.
func checkTrackedSensitiveFiles(projectPath string) []Finding {
	var findings []Finding

	// Check if git is available and this is a git repo
	cmd := exec.Command("git", "ls-files")
	cmd.Dir = projectPath
	output, err := cmd.Output()
	if err != nil {
		return findings // Not a git repo or git not available
	}

	trackedFiles := strings.Split(string(output), "\n")
	for _, tracked := range trackedFiles {
		tracked = strings.TrimSpace(tracked)
		for _, sensitive := range sensitiveFiles {
			basename := filepath.Base(tracked)
			if basename == sensitive || strings.HasPrefix(basename, ".env") {
				findings = append(findings, Finding{
					Category:    CategorySecurity,
					Severity:    SeverityCritical,
					Title:       "Sensitive file tracked by git: " + tracked,
					Description: "This file is being tracked by git and may contain secrets",
					File:        tracked,
					Suggestion:  "Remove from git tracking with 'git rm --cached " + tracked + "' and add to .gitignore",
				})
				break
			}
		}
	}

	return findings
}

// checkSensitiveFiles checks if sensitive files exist in the project.
func checkSensitiveFiles(projectPath string) []Finding {
	var findings []Finding

	for _, sensitive := range sensitiveFiles {
		fullPath := filepath.Join(projectPath, sensitive)
		if _, err := os.Stat(fullPath); err == nil {
			// File exists — check if it's in .gitignore
			gitignorePath := filepath.Join(projectPath, ".gitignore")
			patterns, _ := parseGitignore(gitignorePath)
			if !containsPattern(patterns, sensitive) {
				findings = append(findings, Finding{
					Category:    CategorySecurity,
					Severity:    SeverityError,
					Title:       "Sensitive file not in .gitignore: " + sensitive,
					Description: "The file '" + sensitive + "' exists but is not covered by .gitignore",
					File:        sensitive,
					Suggestion:  "Add '" + sensitive + "' to your .gitignore file",
				})
			}
		}
	}

	return findings
}
