package scanner

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type secretPattern struct {
	Name    string
	Pattern *regexp.Regexp
}

var knownSecretPatterns = []secretPattern{
	{Name: "AWS Access Key ID", Pattern: regexp.MustCompile(`(?i)(AKIA[0-9A-Z]{16})`)},
	{Name: "AWS Secret Key", Pattern: regexp.MustCompile(`(?i)aws_secret_access_key\s*[=:]\s*["']?([A-Za-z0-9/+=]{40})["']?`)},
	{Name: "GitHub Token", Pattern: regexp.MustCompile(`(?i)(ghp_[A-Za-z0-9_]{36,})`)},
	{Name: "Generic API Key", Pattern: regexp.MustCompile(`(?i)(api[_-]?key|apikey|token)\s*[=:]\s*["']([^\s"']{16,})["']`)},
	{Name: "Generic Secret", Pattern: regexp.MustCompile(`(?i)(secret|password|passwd|pwd)\s*[=:]\s*["']([^\s"']{8,})["']`)},
	{Name: "JWT Token", Pattern: regexp.MustCompile(`eyJ[A-Za-z0-9_-]{10,}\.eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_\-]+`)},
	{Name: "Private Key", Pattern: regexp.MustCompile(`-----BEGIN (RSA |EC |DSA |OPENSSH )?PRIVATE KEY-----`)},
	{Name: "Slack Token", Pattern: regexp.MustCompile(`xox[bporas]-[0-9]{10,}-[A-Za-z0-9-]+`)},
	{Name: "Google API Key", Pattern: regexp.MustCompile(`AIza[0-9A-Za-z_-]{35}`)},
	{Name: "Stripe Key", Pattern: regexp.MustCompile(`(?i)(sk|pk)_(live|test)_[0-9a-zA-Z]{24,}`)},
}

var skipExtensions = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".ico": true,
	".svg": true, ".woff": true, ".woff2": true, ".ttf": true, ".eot": true,
	".mp4": true, ".zip": true, ".tar": true, ".gz": true, ".pdf": true,
	".exe": true, ".dll": true, ".so": true, ".lock": true, ".sum": true,
}

var skipDirs = map[string]bool{
	"node_modules": true, ".git": true, "vendor": true, "dist": true,
	"build": true, ".next": true, "__pycache__": true, ".cache": true,
	"coverage": true, "bin": true,
}

// ScanSecrets scans the project for hardcoded secrets.
func ScanSecrets(projectPath string) ([]Finding, error) {
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
		ext := strings.ToLower(filepath.Ext(path))
		if skipExtensions[ext] {
			return nil
		}
		if info.Size() > 1024*1024 {
			return nil
		}
		ff, _ := scanFileForSecrets(path, projectPath)
		findings = append(findings, ff...)
		return nil
	})
	return findings, err
}

func scanFileForSecrets(filePath, projectPath string) ([]Finding, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var findings []Finding
	sc := bufio.NewScanner(file)
	lineNum := 0
	relPath, _ := filepath.Rel(projectPath, filePath)

	for sc.Scan() {
		lineNum++
		line := sc.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") {
			continue
		}
		for _, sp := range knownSecretPatterns {
			if sp.Pattern.MatchString(line) {
				findings = append(findings, Finding{
					Category:    CategorySecurity,
					Severity:    SeverityCritical,
					Title:       "Potential " + sp.Name + " detected",
					Description: "Hardcoded secret matching " + sp.Name + " pattern",
					File:        relPath,
					Line:        lineNum,
					Suggestion:  "Move to environment variable or secrets manager",
				})
				break
			}
		}
		if ef := checkEntropy(line, relPath, lineNum); len(ef) > 0 {
			findings = append(findings, ef...)
		}
	}
	return findings, sc.Err()
}

func checkEntropy(line, file string, lineNum int) []Finding {
	p := regexp.MustCompile(`(?i)(key|token|secret|password|auth)\s*[=:]\s*["']([^\s"']{16,})["']`)
	matches := p.FindStringSubmatch(line)
	if len(matches) < 3 {
		return nil
	}
	entropy := shannonEntropy(matches[2])
	if entropy > 4.5 {
		return []Finding{{
			Category:    CategorySecurity,
			Severity:    SeverityWarning,
			Title:       "High-entropy string detected",
			Description: fmt.Sprintf("Possible secret (entropy: %.2f)", entropy),
			File:        file,
			Line:        lineNum,
			Suggestion:  "Verify this is not a hardcoded secret",
		}}
	}
	return nil
}

func shannonEntropy(s string) float64 {
	if len(s) == 0 {
		return 0
	}
	freq := make(map[rune]float64)
	for _, c := range s {
		freq[c]++
	}
	entropy := 0.0
	length := float64(len(s))
	for _, count := range freq {
		p := count / length
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}
	return entropy
}
