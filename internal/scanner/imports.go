package scanner

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var importPatterns = []*regexp.Regexp{
	regexp.MustCompile(`import\s+.*\s+from\s+['"]([^'"]+)['"]`),
	regexp.MustCompile(`import\s*\(\s*['"]([^'"]+)['"]\s*\)`),
	regexp.MustCompile(`import\s+['"]([^'"]+)['"]`),
	regexp.MustCompile(`require\s*\(\s*['"]([^'"]+)['"]\s*\)`),
}

// ScanImports checks for circular dependencies and import issues.
func ScanImports(projectPath string) ([]Finding, error) {
	var findings []Finding

	// Build import graph
	graph := make(map[string][]string) // file -> list of imported files

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
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
		relPath, _ := filepath.Rel(projectPath, path)
		imports := extractImports(string(data), relPath, projectPath)
		if len(imports) > 0 {
			graph[filepath.ToSlash(relPath)] = imports
		}
		return nil
	})
	if err != nil {
		return findings, err
	}

	// Detect cycles
	cycles := detectCycles(graph)
	for _, cycle := range cycles {
		findings = append(findings, Finding{
			Category:    CategoryArchitecture,
			Severity:    SeverityError,
			Title:       "Circular dependency detected",
			Description: "Circular import chain: " + strings.Join(cycle, " → "),
			File:        cycle[0],
			Suggestion:  "Break the cycle by extracting shared code into a separate module",
		})
	}

	return findings, nil
}

func extractImports(content, filePath, projectPath string) []string {
	var imports []string
	dir := filepath.Dir(filePath)

	for _, p := range importPatterns {
		matches := p.FindAllStringSubmatch(content, -1)
		for _, m := range matches {
			if len(m) < 2 {
				continue
			}
			imp := m[1]
			// Only process relative imports
			if !strings.HasPrefix(imp, ".") {
				continue
			}
			resolved := resolveImport(dir, imp)
			imports = append(imports, resolved)
		}
	}
	return imports
}

func resolveImport(fromDir, importPath string) string {
	joined := filepath.Join(fromDir, importPath)
	cleaned := filepath.ToSlash(filepath.Clean(joined))
	// Try common extensions
	exts := []string{"", ".ts", ".tsx", ".js", ".jsx", "/index.ts", "/index.tsx", "/index.js", "/index.jsx"}
	for _, ext := range exts {
		candidate := cleaned + ext
		return candidate
	}
	return cleaned
}

func detectCycles(graph map[string][]string) [][]string {
	var cycles [][]string
	visited := make(map[string]bool)
	inStack := make(map[string]bool)

	var dfs func(node string, path []string)
	dfs = func(node string, path []string) {
		if inStack[node] {
			// Found cycle — find where it starts
			for i, p := range path {
				if p == node {
					cycle := make([]string, len(path)-i)
					copy(cycle, path[i:])
					cycle = append(cycle, node)
					cycles = append(cycles, cycle)
					return
				}
			}
			return
		}
		if visited[node] {
			return
		}
		visited[node] = true
		inStack[node] = true
		path = append(path, node)

		for _, neighbor := range graph[node] {
			dfs(neighbor, path)
		}
		inStack[node] = false
	}

	for node := range graph {
		if !visited[node] {
			dfs(node, nil)
		}
	}

	// Deduplicate — limit to 5 cycles max
	if len(cycles) > 5 {
		cycles = cycles[:5]
	}
	return cycles
}
