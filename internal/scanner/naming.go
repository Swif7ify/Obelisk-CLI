package scanner

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// NamingRules defines naming conventions for a framework.
type NamingRules struct {
	ComponentPattern *regexp.Regexp // e.g., PascalCase for React
	AssetPattern     *regexp.Regexp // e.g., kebab-case for assets
	ComponentDirs    []string       // directories that should contain components
	AssetDirs        []string       // directories that should contain assets
}

// DefaultJSNamingRules returns naming rules for JS/TS React projects.
func DefaultJSNamingRules() NamingRules {
	return NamingRules{
		ComponentPattern: regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*\.(jsx|tsx)$`),
		AssetPattern:     regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*\.[a-z]+$`),
		ComponentDirs:    []string{"components", "pages", "views", "screens", "containers"},
		AssetDirs:        []string{"assets", "images", "icons", "fonts", "styles"},
	}
}

// ScanNaming checks file/folder naming conventions.
func ScanNaming(projectPath string, rules NamingRules) ([]Finding, error) {
	var findings []Finding
	
	// Create gitignore matcher
	matcher, _ := NewGitignoreMatcher(projectPath)

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		// Skip if path should be ignored
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
		dir := filepath.Dir(relPath)
		name := info.Name()
		ext := filepath.Ext(name)

		// Check component files (JSX/TSX should be PascalCase)
		if ext == ".jsx" || ext == ".tsx" {
			if isInDirs(dir, rules.ComponentDirs) || containsComponentDir(dir) {
				if !rules.ComponentPattern.MatchString(name) && !strings.HasPrefix(name, "index") {
					findings = append(findings, Finding{
						Category:    CategoryNaming,
						Severity:    SeverityWarning,
						Title:       "Component naming violation",
						Description: "React components should use PascalCase",
						File:        relPath,
						Suggestion:  "Rename to " + toPascalCase(strings.TrimSuffix(name, ext)) + ext,
					})
				}
			}
		}

		// Check for spaces in filenames
		if strings.Contains(name, " ") {
			findings = append(findings, Finding{
				Category:    CategoryNaming,
				Severity:    SeverityWarning,
				Title:       "Filename contains spaces",
				Description: "Filenames should not contain spaces",
				File:        relPath,
				Suggestion:  "Use kebab-case or camelCase instead",
			})
		}

		// Check for uppercase in non-component dirs
		if isInDirs(dir, rules.AssetDirs) {
			if hasUpperCase(strings.TrimSuffix(name, ext)) && ext != ".jsx" && ext != ".tsx" {
				findings = append(findings, Finding{
					Category:    CategoryNaming,
					Severity:    SeverityInfo,
					Title:       "Asset naming convention",
					Description: "Asset files should use kebab-case",
					File:        relPath,
					Suggestion:  "Rename to lowercase kebab-case",
				})
			}
		}

		return nil
	})

	return findings, err
}

func isInDirs(dir string, targets []string) bool {
	parts := strings.Split(filepath.ToSlash(dir), "/")
	for _, p := range parts {
		for _, t := range targets {
			if strings.EqualFold(p, t) {
				return true
			}
		}
	}
	return false
}

func containsComponentDir(dir string) bool {
	lower := strings.ToLower(dir)
	return strings.Contains(lower, "component") || strings.Contains(lower, "page") ||
		strings.Contains(lower, "view") || strings.Contains(lower, "screen")
}

func hasUpperCase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func toPascalCase(s string) string {
	parts := regexp.MustCompile(`[-_\s]+`).Split(s, -1)
	var result strings.Builder
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		result.WriteString(strings.ToUpper(p[:1]) + p[1:])
	}
	return result.String()
}
