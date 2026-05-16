package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var techDebtRegex = regexp.MustCompile(`(?i)\b(TODO|FIXME|HACK|XXX)\b`)

// ScanTechDebt hunts for unresolved technical debt comments (TODO, FIXME, etc.)
func ScanTechDebt(projectPath string) ([]Finding, error) {
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
		
		// Only check tech debt for source code files
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

		debtCount := 0

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		sc := bufio.NewScanner(file)
		for sc.Scan() {
			line := sc.Text()
			if techDebtRegex.MatchString(line) {
				debtCount++
			}
		}

		relPath, _ := filepath.Rel(projectPath, path)
		if relPath == "" {
			relPath = path
		}

		if debtCount > 10 {
			findings = append(findings, Finding{
				Category:    CategoryQuality,
				Severity:    SeverityError,
				Title:       "High Technical Debt Accumulation",
				Description: fmt.Sprintf("File contains %d unresolved TODO/FIXME/HACK comments.", debtCount),
				File:        relPath,
				Suggestion:  "Resolve lingering tech debt before it becomes unmanageable.",
			})
		} else if debtCount > 3 {
			findings = append(findings, Finding{
				Category:    CategoryQuality,
				Severity:    SeverityWarning,
				Title:       "Lingering Technical Debt",
				Description: fmt.Sprintf("File contains %d TODO/FIXME/HACK comments.", debtCount),
				File:        relPath,
				Suggestion:  "Review and address these comments to maintain code quality.",
			})
		}

		return nil
	})

	return findings, err
}
