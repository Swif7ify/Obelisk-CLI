package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanSecrets(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()

	tests := []struct {
		name          string
		files         map[string]string
		wantFindings  int
		wantSeverity  string
		wantContains  string
	}{
		{
			name: "detects hardcoded API key",
			files: map[string]string{
				"config.js": `const API_KEY = "sk_live_1234567890abcdef";`,
			},
			wantFindings: 1,
			wantSeverity: "critical",
			wantContains: "API key",
		},
		{
			name: "detects AWS credentials",
			files: map[string]string{
				"aws.js": `const AWS_ACCESS_KEY_ID = "AKIAIOSFODNN7EXAMPLE";`,
			},
			wantFindings: 1,
			wantSeverity: "critical",
			wantContains: "AWS",
		},
		{
			name: "detects JWT token",
			files: map[string]string{
				"auth.js": `const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U";`,
			},
			wantFindings: 1,
			wantSeverity: "critical",
			wantContains: "JWT",
		},
		{
			name: "detects private key",
			files: map[string]string{
				"key.pem": `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA1234567890
-----END RSA PRIVATE KEY-----`,
			},
			wantFindings: 1,
			wantSeverity: "critical",
			wantContains: "Private key",
		},
		{
			name: "ignores safe patterns",
			files: map[string]string{
				"example.js": `const API_KEY = process.env.API_KEY;`,
			},
			wantFindings: 0,
		},
		{
			name: "ignores .env.example",
			files: map[string]string{
				".env.example": `API_KEY=your_key_here`,
			},
			wantFindings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test files
			testDir := filepath.Join(tmpDir, tt.name)
			if err := os.MkdirAll(testDir, 0755); err != nil {
				t.Fatal(err)
			}

			for filename, content := range tt.files {
				filePath := filepath.Join(testDir, filename)
				if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
			}

			// Run scan
			findings, err := ScanSecrets(testDir)
			if err != nil {
				t.Fatalf("ScanSecrets() error = %v", err)
			}

			// Check findings count
			if len(findings) != tt.wantFindings {
				t.Errorf("ScanSecrets() got %d findings, want %d", len(findings), tt.wantFindings)
			}

			// Check severity and message if findings expected
			if tt.wantFindings > 0 && len(findings) > 0 {
				if findings[0].Severity != tt.wantSeverity {
					t.Errorf("ScanSecrets() severity = %v, want %v", findings[0].Severity, tt.wantSeverity)
				}
				if tt.wantContains != "" && !contains(findings[0].Message, tt.wantContains) {
					t.Errorf("ScanSecrets() message = %v, want to contain %v", findings[0].Message, tt.wantContains)
				}
			}
		})
	}
}

func TestCalculateEntropy(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantMin float64
	}{
		{
			name:    "high entropy random string",
			input:   "aB3$xY9#mK2@pL5",
			wantMin: 3.0,
		},
		{
			name:    "low entropy repeated chars",
			input:   "aaaaaaaaaa",
			wantMin: 0.0,
		},
		{
			name:    "medium entropy",
			input:   "password123",
			wantMin: 2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entropy := calculateEntropy(tt.input)
			if entropy < tt.wantMin {
				t.Errorf("calculateEntropy() = %v, want >= %v", entropy, tt.wantMin)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Made with Bob
