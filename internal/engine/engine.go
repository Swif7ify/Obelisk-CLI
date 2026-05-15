package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Swif7ify/Obelisk-CLI/adapters"
	"github.com/Swif7ify/Obelisk-CLI/internal/ai"
	"github.com/Swif7ify/Obelisk-CLI/internal/detector"
	"github.com/Swif7ify/Obelisk-CLI/internal/linter"
	"github.com/Swif7ify/Obelisk-CLI/internal/scanner"
)

// Config holds the engine configuration.
type Config struct {
	ProjectPath string
	APIKey      string
	Model       string
	SkipAI      bool
	Verbose     bool
}

// Result holds the complete engine output.
type Result struct {
	ScanResult *scanner.ScanResult
	Report     *ai.HealthReport
	Detection  detector.DetectionResult
}

// OnPhaseChange is a callback for scan phase updates.
type OnPhaseChange func(phase string)

// Run executes the full scan pipeline.
func Run(cfg Config, onPhase OnPhaseChange) (*Result, error) {
	projectPath := cfg.ProjectPath
	if projectPath == "" {
		var err error
		projectPath, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}
	}

	// Resolve absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, fmt.Errorf("invalid project path: %w", err)
	}

	// Step 1: Detect project type
	if onPhase != nil {
		onPhase("Detecting project type...")
	}
	detection := detector.Detect(absPath)

	// Step 2: Load adapter
	adapter := getAdapter(detection.Type)

	// Step 3: Count files/dirs
	fileCount, dirCount := countFilesAndDirs(absPath)

	result := &scanner.ScanResult{
		ProjectPath: absPath,
		ProjectType: string(detection.Type),
		FileCount:   fileCount,
		DirCount:    dirCount,
	}

	// Step 4: Run all scanners
	if onPhase != nil {
		onPhase("Scanning for secrets...")
	}
	secretFindings, _ := scanner.ScanSecrets(absPath)
	result.Findings = append(result.Findings, secretFindings...)

	if onPhase != nil {
		onPhase("Validating .gitignore...")
	}
	gitignoreFindings, _ := scanner.ScanGitignore(absPath)
	result.Findings = append(result.Findings, gitignoreFindings...)

	if onPhase != nil {
		onPhase("Auditing dependencies...")
	}
	depFindings, _ := scanner.ScanDependencies(absPath)
	result.Findings = append(result.Findings, depFindings...)

	if onPhase != nil {
		onPhase("Checking naming conventions...")
	}
	if adapter != nil {
		namingFindings, _ := scanner.ScanNaming(absPath, adapter.NamingRules())
		result.Findings = append(result.Findings, namingFindings...)
	}

	if onPhase != nil {
		onPhase("Analyzing imports...")
	}
	importFindings, _ := scanner.ScanImports(absPath)
	result.Findings = append(result.Findings, importFindings...)

	// Step 5: Orchestrated Linting (ESLint)
	if onPhase != nil {
		onPhase("Running linters...")
	}
	lintFindings, _ := linter.RunESLint(absPath)
	result.Findings = append(result.Findings, lintFindings...)

	// Step 6: Build directory tree for AI Vibe Check
	dirTree := buildDirTree(absPath, "", 0, 3)
	result.DirTree = dirTree

	// Step 7: AI analysis
	var report *ai.HealthReport
	if !cfg.SkipAI {
		if onPhase != nil {
			onPhase("Consulting AI brain...")
		}
		aiClient, err := ai.NewClient(cfg.APIKey, cfg.Model)
		if err != nil {
			// Fallback to non-AI report
			report = ai.FallbackReport(
				result.CountBySeverity(scanner.SeverityCritical),
				result.CountBySeverity(scanner.SeverityError),
				result.CountBySeverity(scanner.SeverityWarning),
				result.CountBySeverity(scanner.SeverityInfo),
			)
		} else {
			prompt := ai.BuildPrompt(result)
			ctx := context.Background()
			response, err := aiClient.GenerateContent(ctx, prompt)
			if err != nil {
				report = ai.FallbackReport(
					result.CountBySeverity(scanner.SeverityCritical),
					result.CountBySeverity(scanner.SeverityError),
					result.CountBySeverity(scanner.SeverityWarning),
					result.CountBySeverity(scanner.SeverityInfo),
				)
			} else {
				report, err = ai.ParseReport(response)
				if err != nil {
					report = ai.FallbackReport(
						result.CountBySeverity(scanner.SeverityCritical),
						result.CountBySeverity(scanner.SeverityError),
						result.CountBySeverity(scanner.SeverityWarning),
						result.CountBySeverity(scanner.SeverityInfo),
					)
				}
			}
		}
	} else {
		report = ai.FallbackReport(
			result.CountBySeverity(scanner.SeverityCritical),
			result.CountBySeverity(scanner.SeverityError),
			result.CountBySeverity(scanner.SeverityWarning),
			result.CountBySeverity(scanner.SeverityInfo),
		)
	}

	return &Result{
		ScanResult: result,
		Report:     report,
		Detection:  detection,
	}, nil
}

func getAdapter(t detector.ProjectType) adapters.Adapter {
	switch t {
	case detector.TypeJavaScript, detector.TypeTypeScript, detector.TypeReact, detector.TypeNextJS:
		return &adapters.JavaScriptAdapter{}
	case detector.TypeLaravel:
		return &adapters.LaravelAdapter{}
	case detector.TypeGolang:
		return &adapters.GolangAdapter{}
	default:
		return &adapters.JavaScriptAdapter{} // Default fallback
	}
}

func countFilesAndDirs(root string) (int, int) {
	files, dirs := 0, 0
	
	// Create gitignore matcher
	matcher, _ := scanner.NewGitignoreMatcher(root)
	
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
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
			dirs++
		} else {
			files++
		}
		return nil
	})
	return files, dirs
}

// buildDirTree generates a text representation of the project directory tree.
func buildDirTree(root, prefix string, depth, maxDepth int) string {
	if depth >= maxDepth {
		return ""
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return ""
	}

	var sb strings.Builder
	skipNames := map[string]bool{
		"node_modules": true, ".git": true, "vendor": true, "dist": true,
		"build": true, ".next": true, "__pycache__": true, ".cache": true,
		"coverage": true, "bin": true, ".github": true,
	}

	// Filter visible entries
	var visible []os.DirEntry
	for _, e := range entries {
		if !skipNames[e.Name()] {
			visible = append(visible, e)
		}
	}

	for i, entry := range visible {
		isLast := i == len(visible)-1
		connector := "├── "
		childPrefix := prefix + "│   "
		if isLast {
			connector = "└── "
			childPrefix = prefix + "    "
		}

		if entry.IsDir() {
			sb.WriteString(prefix + connector + entry.Name() + "/\n")
			sb.WriteString(buildDirTree(filepath.Join(root, entry.Name()), childPrefix, depth+1, maxDepth))
		} else {
			sb.WriteString(prefix + connector + entry.Name() + "\n")
		}
	}

	return sb.String()
}
