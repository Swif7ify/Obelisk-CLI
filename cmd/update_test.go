package cmd

import (
	"testing"
)

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		name    string
		current string
		latest  string
		want    bool
	}{
		{
			name:    "current is newer major",
			current: "2.0.0",
			latest:  "1.0.0",
			want:    true,
		},
		{
			name:    "current is newer minor",
			current: "1.5.0",
			latest:  "1.4.0",
			want:    true,
		},
		{
			name:    "current is newer patch",
			current: "1.0.5",
			latest:  "1.0.4",
			want:    true,
		},
		{
			name:    "current is older",
			current: "1.0.0",
			latest:  "2.0.0",
			want:    false,
		},
		{
			name:    "versions are equal",
			current: "1.0.0",
			latest:  "1.0.0",
			want:    false,
		},
		{
			name:    "current has more parts",
			current: "1.0.0.1",
			latest:  "1.0.0",
			want:    true,
		},
		{
			name:    "latest has more parts",
			current: "1.0.0",
			latest:  "1.0.0.1",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isNewerVersion(tt.current, tt.latest)
			if got != tt.want {
				t.Errorf("isNewerVersion(%q, %q) = %v, want %v", tt.current, tt.latest, got, tt.want)
			}
		})
	}
}

func TestGetAssetURL(t *testing.T) {
	release := &GitHubRelease{
		HTMLURL: "https://github.com/test/repo/releases/tag/v1.0.0",
		Assets: []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		}{
			{
				Name:               "obelisk-windows-amd64.exe",
				BrowserDownloadURL: "https://github.com/test/repo/releases/download/v1.0.0/obelisk-windows-amd64.exe",
			},
			{
				Name:               "obelisk-linux-amd64",
				BrowserDownloadURL: "https://github.com/test/repo/releases/download/v1.0.0/obelisk-linux-amd64",
			},
			{
				Name:               "ObeliskCLI-1.0.0-x64.msi",
				BrowserDownloadURL: "https://github.com/test/repo/releases/download/v1.0.0/ObeliskCLI-1.0.0-x64.msi",
			},
		},
	}

	tests := []struct {
		name    string
		keyword string
		want    string
	}{
		{
			name:    "find windows exe",
			keyword: "exe",
			want:    "https://github.com/test/repo/releases/download/v1.0.0/obelisk-windows-amd64.exe",
		},
		{
			name:    "find linux binary",
			keyword: "linux",
			want:    "https://github.com/test/repo/releases/download/v1.0.0/obelisk-linux-amd64",
		},
		{
			name:    "find msi installer",
			keyword: "msi",
			want:    "https://github.com/test/repo/releases/download/v1.0.0/ObeliskCLI-1.0.0-x64.msi",
		},
		{
			name:    "keyword not found returns release page",
			keyword: "notfound",
			want:    "https://github.com/test/repo/releases/tag/v1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAssetURL(release, tt.keyword)
			if got != tt.want {
				t.Errorf("getAssetURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersionComparison(t *testing.T) {
	// Test version comparison edge cases
	tests := []struct {
		name    string
		v1      string
		v2      string
		wantNewer bool
	}{
		{"1.0.0 vs 0.9.9", "1.0.0", "0.9.9", true},
		{"0.9.9 vs 1.0.0", "0.9.9", "1.0.0", false},
		{"1.10.0 vs 1.9.0", "1.10.0", "1.9.0", true},
		{"1.0.10 vs 1.0.9", "1.0.10", "1.0.9", true},
		{"2.0.0 vs 1.99.99", "2.0.0", "1.99.99", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isNewerVersion(tt.v1, tt.v2)
			if got != tt.wantNewer {
				t.Errorf("isNewerVersion(%q, %q) = %v, want %v", tt.v1, tt.v2, got, tt.wantNewer)
			}
		})
	}
}

// Made with Bob
