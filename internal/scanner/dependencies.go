package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// PackageJSON represents relevant fields from package.json.
type PackageJSON struct {
	Name            string            `json:"name"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

var knownDeprecatedPackages = map[string]string{
	"request":          "Deprecated since 2020. Use 'node-fetch', 'axios', or 'got'",
	"tslint":           "Deprecated. Migrate to ESLint with @typescript-eslint",
	"node-sass":        "Deprecated. Use 'sass' (Dart Sass) instead",
	"moment":           "Deprecated. Use 'date-fns', 'dayjs', or 'luxon'",
	"enzyme":           "No longer maintained. Use React Testing Library",
	"create-react-app": "Deprecated. Use Vite, Next.js, or Remix",
}

var knownVulnerablePackages = map[string]string{
	"event-stream":   "Compromised in 2018 (malicious code injection)",
	"flatmap-stream": "Malicious package",
	"colors":         "Author sabotaged v1.4.1+ — pin to 1.4.0",
	"faker":          "Author sabotaged v6.6.6+ — use @faker-js/faker",
	"node-ipc":       "Protestware in v10.1.1+",
}

// ScanDependencies analyzes package.json for deprecated and vulnerable packages.
func ScanDependencies(projectPath string) ([]Finding, error) {
	var findings []Finding
	pkgPath := filepath.Join(projectPath, "package.json")
	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		return findings, nil
	}
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return findings, err
	}
	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		findings = append(findings, Finding{
			Category: CategoryQuality, Severity: SeverityError,
			Title: "Invalid package.json", Description: "Failed to parse: " + err.Error(),
			File: "package.json", Suggestion: "Fix the JSON syntax",
		})
		return findings, nil
	}
	allDeps := make(map[string]string)
	for k, v := range pkg.Dependencies {
		allDeps[k] = v
	}
	for k, v := range pkg.DevDependencies {
		allDeps[k] = v
	}
	for dep := range allDeps {
		if reason, ok := knownDeprecatedPackages[dep]; ok {
			findings = append(findings, Finding{
				Category: CategoryDependency, Severity: SeverityWarning,
				Title: "Deprecated package: " + dep, Description: reason,
				File: "package.json", Suggestion: "Replace with a modern alternative",
			})
		}
		if reason, ok := knownVulnerablePackages[dep]; ok {
			findings = append(findings, Finding{
				Category: CategorySecurity, Severity: SeverityCritical,
				Title: "Vulnerable package: " + dep, Description: reason,
				File: "package.json", Suggestion: "Remove or replace immediately",
			})
		}
	}
	findings = append(findings, checkUnusedDeps(projectPath, pkg.Dependencies)...)
	return findings, nil
}

func checkUnusedDeps(projectPath string, deps map[string]string) []Finding {
	var findings []Finding
	for dep := range deps {
		if !isDepUsed(projectPath, dep) {
			findings = append(findings, Finding{
				Category: CategoryDependency, Severity: SeverityInfo,
				Title: "Potentially unused: " + dep,
				Description: "No import found in source files",
				File: "package.json", Suggestion: "Remove with 'npm uninstall " + dep + "'",
			})
		}
	}
	return findings
}

func isDepUsed(projectPath, dep string) bool {
	found := false
	_ = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || found {
			return filepath.SkipDir
		}
		if info.IsDir() {
			n := info.Name()
			if n == "node_modules" || n == ".git" || n == "dist" || n == "build" || n == ".next" {
				return filepath.SkipDir
			}
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".js" && ext != ".jsx" && ext != ".ts" && ext != ".tsx" {
			return nil
		}
		data, e := os.ReadFile(path)
		if e != nil {
			return nil
		}
		c := string(data)
		if strings.Contains(c, `"`+dep+`"`) || strings.Contains(c, `'`+dep+`'`) ||
			strings.Contains(c, `"`+dep+"/") || strings.Contains(c, `'`+dep+"/") {
			found = true
			return filepath.SkipDir
		}
		return nil
	})
	return found
}
