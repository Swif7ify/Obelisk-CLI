package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanGitignore(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name         string
		gitignore    string
		envFiles     []string
		wantFindings int
		wantSeverity string
	}{
		{
			name: "missing .gitignore",
			wantFindings: 1,
			wantSeverity: "error",
		},
		{
			name: "valid .gitignore with node_modules",
			gitignore: `node_modules/
.env
*.log`,
			wantFindings: 0,
		},
		{
			name: "missing node_modules in .gitignore",
			gitignore: `.env
*.log`,
			wantFindings: 1,
			wantSeverity: "warning",
		},
		{
			name: "missing .env in .gitignore",
			gitignore: `node_modules/
*.log`,
			wantFindings: 1,
			wantSeverity: "critical",
		},
		{
			name: ".env file not ignored",
			gitignore: `node_modules/
*.log`,
			envFiles: []string{".env"},
			wantFindings: 2, // Missing .env pattern + tracked .env file
			wantSeverity: "critical",
		},
		{
			name: ".env.example is safe",
			gitignore: `node_modules/
.env
*.log`,
			envFiles: []string{".env.example"},
			wantFindings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test directory
			testDir := filepath.Join(tmpDir, tt.name)
			if err := os.MkdirAll(testDir, 0755); err != nil {
				t.Fatal(err)
			}

			// Create .gitignore if specified
			if tt.gitignore != "" {
				gitignorePath := filepath.Join(testDir, ".gitignore")
				if err := os.WriteFile(gitignorePath, []byte(tt.gitignore), 0644); err != nil {
					t.Fatal(err)
				}
			}

			// Create env files if specified
			for _, envFile := range tt.envFiles {
				envPath := filepath.Join(testDir, envFile)
				if err := os.WriteFile(envPath, []byte("SECRET=value"), 0644); err != nil {
					t.Fatal(err)
				}
			}

			// Initialize git repo (required for some checks)
			gitDir := filepath.Join(testDir, ".git")
			if err := os.MkdirAll(gitDir, 0755); err != nil {
				t.Fatal(err)
			}

			// Run scan
			findings, err := ScanGitignore(testDir)
			if err != nil {
				t.Fatalf("ScanGitignore() error = %v", err)
			}

			// Check findings count
			if len(findings) != tt.wantFindings {
				t.Errorf("ScanGitignore() got %d findings, want %d", len(findings), tt.wantFindings)
				for _, f := range findings {
					t.Logf("Finding: %s - %s (severity: %s)", f.File, f.Message, f.Severity)
				}
			}

			// Check severity if findings expected
			if tt.wantFindings > 0 && len(findings) > 0 && tt.wantSeverity != "" {
				found := false
				for _, finding := range findings {
					if finding.Severity == tt.wantSeverity {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("ScanGitignore() no finding with severity %v", tt.wantSeverity)
				}
			}
		})
	}
}

func TestParseGitignore(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name: "simple patterns",
			content: `node_modules/
.env
*.log`,
			want: []string{"node_modules/", ".env", "*.log"},
		},
		{
			name: "with comments and empty lines",
			content: `# Dependencies
node_modules/

# Environment
.env

# Logs
*.log`,
			want: []string{"node_modules/", ".env", "*.log"},
		},
		{
			name: "negation patterns",
			content: `*.log
!important.log`,
			want: []string{"*.log", "!important.log"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tmpDir, tt.name+".gitignore")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			// Parse
			got, err := parseGitignore(testFile)
			if err != nil {
				t.Fatalf("parseGitignore() error = %v", err)
			}

			// Check patterns
			if len(got) != len(tt.want) {
				t.Errorf("parseGitignore() got %d patterns, want %d", len(got), len(tt.want))
			}

			for i, pattern := range tt.want {
				if i >= len(got) {
					break
				}
				if got[i] != pattern {
					t.Errorf("parseGitignore() pattern[%d] = %q, want %q", i, got[i], pattern)
				}
			}
		})
	}
}

func TestCheckSensitiveFiles(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name         string
		files        []string
		wantFindings int
	}{
		{
			name:         "no sensitive files",
			files:        []string{"README.md", "src/index.js"},
			wantFindings: 0,
		},
		{
			name:         "has .env file",
			files:        []string{".env", "src/index.js"},
			wantFindings: 1,
		},
		{
			name:         ".env.example is safe",
			files:        []string{".env.example", "src/index.js"},
			wantFindings: 0,
		},
		{
			name:         "multiple sensitive files",
			files:        []string{".env", "config.json", "secrets.yml"},
			wantFindings: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test directory
			testDir := filepath.Join(tmpDir, tt.name)
			if err := os.MkdirAll(testDir, 0755); err != nil {
				t.Fatal(err)
			}

			// Create test files
			for _, file := range tt.files {
				filePath := filepath.Join(testDir, file)
				if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
					t.Fatal(err)
				}
			}

			// Run check
			findings := checkSensitiveFiles(testDir)

			// Verify findings count
			if len(findings) != tt.wantFindings {
				t.Errorf("checkSensitiveFiles() got %d findings, want %d", len(findings), tt.wantFindings)
			}
		})
	}
}

// Made with Bob
