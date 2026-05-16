package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanNaming(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name         string
		files        []string
		rules        NamingRules
		wantFindings int
		wantContains string
	}{
		{
			name: "React component should be PascalCase",
			files: []string{
				"src/components/userProfile.jsx",
			},
			rules: NamingRules{
				ComponentCase: "PascalCase",
				FileCase:      "kebab-case",
			},
			wantFindings: 1,
			wantContains: "PascalCase",
		},
		{
			name: "Valid React component naming",
			files: []string{
				"src/components/UserProfile.jsx",
			},
			rules: NamingRules{
				ComponentCase: "PascalCase",
				FileCase:      "kebab-case",
			},
			wantFindings: 0,
		},
		{
			name: "Asset files should be kebab-case",
			files: []string{
				"assets/UserIcon.png",
			},
			rules: NamingRules{
				ComponentCase: "PascalCase",
				FileCase:      "kebab-case",
			},
			wantFindings: 1,
			wantContains: "kebab-case",
		},
		{
			name: "Valid asset naming",
			files: []string{
				"assets/user-icon.png",
			},
			rules: NamingRules{
				ComponentCase: "PascalCase",
				FileCase:      "kebab-case",
			},
			wantFindings: 0,
		},
		{
			name: "Multiple violations",
			files: []string{
				"src/components/userProfile.jsx",
				"assets/UserIcon.png",
				"utils/helperFunctions.js",
			},
			rules: NamingRules{
				ComponentCase: "PascalCase",
				FileCase:      "kebab-case",
			},
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
				if err := os.WriteFile(filePath, []byte("// test content"), 0644); err != nil {
					t.Fatal(err)
				}
			}

			// Run scan
			findings, err := ScanNaming(testDir, tt.rules)
			if err != nil {
				t.Fatalf("ScanNaming() error = %v", err)
			}

			// Check findings count
			if len(findings) != tt.wantFindings {
				t.Errorf("ScanNaming() got %d findings, want %d", len(findings), tt.wantFindings)
				for _, f := range findings {
					t.Logf("Finding: %s - %s", f.File, f.Message)
				}
			}

			// Check message content if findings expected
			if tt.wantFindings > 0 && len(findings) > 0 && tt.wantContains != "" {
				found := false
				for _, finding := range findings {
					if contains(finding.Message, tt.wantContains) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("ScanNaming() no finding contains %v", tt.wantContains)
				}
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"userProfile", "UserProfile"},
		{"user-profile", "UserProfile"},
		{"user_profile", "UserProfile"},
		{"UserProfile", "UserProfile"},
		{"user", "User"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toPascalCase(tt.input)
			if got != tt.want {
				t.Errorf("toPascalCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNamingRulesValidation(t *testing.T) {
	tests := []struct {
		name  string
		rules NamingRules
		valid bool
	}{
		{
			name: "valid rules",
			rules: NamingRules{
				ComponentCase: "PascalCase",
				FileCase:      "kebab-case",
			},
			valid: true,
		},
		{
			name: "empty rules",
			rules: NamingRules{
				ComponentCase: "",
				FileCase:      "",
			},
			valid: true, // Should handle gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			_, err := ScanNaming(tmpDir, tt.rules)
			if tt.valid && err != nil {
				t.Errorf("ScanNaming() unexpected error = %v", err)
			}
		})
	}
}

// Made with Bob
