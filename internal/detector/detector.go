package detector

import (
	"os"
	"path/filepath"
)

// ProjectType represents a detected framework/language.
type ProjectType string

const (
	TypeJavaScript ProjectType = "javascript"
	TypeTypeScript ProjectType = "typescript"
	TypeReact      ProjectType = "react"
	TypeNextJS     ProjectType = "nextjs"
	TypeLaravel    ProjectType = "laravel"
	TypeGolang     ProjectType = "golang"
	TypePython     ProjectType = "python"
	TypeUnknown    ProjectType = "unknown"
)

// DetectionResult holds the detected project information.
type DetectionResult struct {
	Type       ProjectType
	Framework  string
	ConfigFile string
}

// Detect scans the project root and identifies the framework/language.
func Detect(projectPath string) DetectionResult {
	// Check for Next.js (has next.config.*)
	nextConfigs := []string{"next.config.js", "next.config.mjs", "next.config.ts"}
	for _, nc := range nextConfigs {
		if fileExists(filepath.Join(projectPath, nc)) {
			return DetectionResult{Type: TypeNextJS, Framework: "Next.js", ConfigFile: nc}
		}
	}

	// Check for package.json (JS/TS ecosystem)
	pkgPath := filepath.Join(projectPath, "package.json")
	if fileExists(pkgPath) {
		// Check for React indicators
		if fileExists(filepath.Join(projectPath, "src", "App.tsx")) ||
			fileExists(filepath.Join(projectPath, "src", "App.jsx")) ||
			fileExists(filepath.Join(projectPath, "src", "App.js")) {
			return DetectionResult{Type: TypeReact, Framework: "React", ConfigFile: "package.json"}
		}

		// Check for tsconfig
		if fileExists(filepath.Join(projectPath, "tsconfig.json")) {
			return DetectionResult{Type: TypeTypeScript, Framework: "TypeScript", ConfigFile: "tsconfig.json"}
		}

		return DetectionResult{Type: TypeJavaScript, Framework: "JavaScript", ConfigFile: "package.json"}
	}

	// Check for Laravel (composer.json + artisan)
	if fileExists(filepath.Join(projectPath, "composer.json")) && fileExists(filepath.Join(projectPath, "artisan")) {
		return DetectionResult{Type: TypeLaravel, Framework: "Laravel", ConfigFile: "composer.json"}
	}

	// Check for Go
	if fileExists(filepath.Join(projectPath, "go.mod")) {
		return DetectionResult{Type: TypeGolang, Framework: "Go", ConfigFile: "go.mod"}
	}

	// Check for Python
	if fileExists(filepath.Join(projectPath, "requirements.txt")) || fileExists(filepath.Join(projectPath, "pyproject.toml")) {
		return DetectionResult{Type: TypePython, Framework: "Python", ConfigFile: "requirements.txt"}
	}

	return DetectionResult{Type: TypeUnknown, Framework: "Unknown", ConfigFile: ""}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
