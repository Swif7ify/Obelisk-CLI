package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var branchRegex = regexp.MustCompile(`\b(if|else if|for|while|case|catch)\b|(&&|\|\||\?)`)

// ScanComplexity calculates a rudimentary Cyclomatic Complexity score for source files.
// It flags files that have excessive branching, indicating "Spaghetti Code".
func ScanComplexity(projectPath string) ([]Finding, error) {
	var findings []Finding

	matcher, _ := NewGitignoreMatcher(projectPath)

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if matcher.ShouldIgnore(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		
		// Only check complexity for source code files
		validExts := map[string]bool{
			".js": true, ".ts": true, ".jsx": true, ".tsx": true,
			".go": true, ".php": true, ".py": true, ".rb": true,
			".java": true, ".c": true, ".cpp": true, ".cs": true,
		}
		if !validExts[ext] {
			return nil
		}

		if info.Size() > 1024*1024 { // Skip files > 1MB
			return nil
		}

		complexityScore := 1 // Base complexity is 1

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		sc := bufio.NewScanner(file)
		for sc.Scan() {
			line := sc.Text()
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "/*") {
				continue
			}
			matches := branchRegex.FindAllString(line, -1)
			complexityScore += len(matches)
		}

		relPath, _ := filepath.Rel(projectPath, path)
		if relPath == "" {
			relPath = path
		}

		if complexityScore > 50 {
			findings = append(findings, Finding{
				Category:    CategoryArchitecture,
				Severity:    SeverityError,
				Title:       "Severe Spaghetti Code Detected",
				Description: fmt.Sprintf("File has a Cyclomatic Complexity score of %d (>50). It is highly unmaintainable.", complexityScore),
				File:        relPath,
				Suggestion:  "Refactor the file into smaller, modular functions or components.",
			})
		} else if complexityScore > 25 {
			findings = append(findings, Finding{
				Category:    CategoryArchitecture,
				Severity:    SeverityWarning,
				Title:       "High Cyclomatic Complexity",
				Description: fmt.Sprintf("File has a Cyclomatic Complexity score of %d (>25). Logic is hard to follow.", complexityScore),
				File:        relPath,
				Suggestion:  "Consider breaking down large functions to reduce branch density.",
			})
		}

		return nil
	})

	return findings, err
}
