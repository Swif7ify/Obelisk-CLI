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

		relPath, _ := filepath.Rel(projectPath, path)
		if relPath == "" {
			relPath = path
		}

		baseUpper := strings.ToUpper(filepath.Base(path))
		isDebtFile := false
		if baseUpper == "TODO" || baseUpper == "TODO.TXT" || baseUpper == "TODO.MD" ||
			baseUpper == "FIXME" || baseUpper == "FIXME.TXT" || baseUpper == "FIXME.MD" ||
			baseUpper == "HACK" || baseUpper == "HACK.TXT" || baseUpper == "HACK.MD" {
			isDebtFile = true
			findings = append(findings, Finding{
				Category:    CategoryQuality,
				Severity:    SeverityWarning,
				Title:       "Dedicated Technical Debt File",
				Description: fmt.Sprintf("Found a dedicated technical debt file: %s", filepath.Base(path)),
				File:        relPath,
				Suggestion:  "Resolve the items in this file and remove it from the repository.",
			})
		}

		ext := strings.ToLower(filepath.Ext(path))
		
		// Check tech debt for source code files, docs, config
		validExts := map[string]bool{
			".js": true, ".ts": true, ".jsx": true, ".tsx": true,
			".go": true, ".php": true, ".py": true, ".rb": true,
			".java": true, ".c": true, ".cpp": true, ".cs": true,
			".json": true, ".md": true, ".txt": true, "": true,
		}
		if !validExts[ext] && !isDebtFile {
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

		if debtCount > 10 {
			findings = append(findings, Finding{
				Category:    CategoryQuality,
				Severity:    SeverityWarning,
				Title:       "High Technical Debt Accumulation",
				Description: fmt.Sprintf("File contains %d unresolved TODO/FIXME/HACK comment(s).", debtCount),
				File:        relPath,
				Suggestion:  "Resolve lingering tech debt before it becomes unmanageable.",
			})
		} else if debtCount > 3 {
			findings = append(findings, Finding{
				Category:    CategoryQuality,
				Severity:    SeverityWarning,
				Title:       "Lingering Technical Debt",
				Description: fmt.Sprintf("File contains %d TODO/FIXME/HACK comment(s).", debtCount),
				File:        relPath,
				Suggestion:  "Review and address these comments to maintain code quality.",
			})
		} else if debtCount > 0 {
			findings = append(findings, Finding{
				Category:    CategoryQuality,
				Severity:    SeverityInfo,
				Title:       "Technical Debt Comment Detected",
				Description: fmt.Sprintf("File contains %d TODO/FIXME/HACK comment(s).", debtCount),
				File:        relPath,
				Suggestion:  "Address the comment to keep the codebase clean.",
			})
		}

		return nil
	})

	return findings, err
}
