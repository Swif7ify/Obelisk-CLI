package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	PublishedAt time.Time `json:"published_at"`
	HTMLURL     string    `json:"html_url"`
	Assets      []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for updates",
	Long: `Check if a newer version of Obelisk CLI is available.
	
This command checks GitHub releases for the latest version and compares it
with your currently installed version.`,
	RunE: runUpdate,
}

var (
	updateCheckOnly bool
	updateAutoYes   bool
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVar(&updateCheckOnly, "check", false, "Only check for updates, don't prompt to install")
	updateCmd.Flags().BoolVarP(&updateAutoYes, "yes", "y", false, "Automatically answer yes to update prompt")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	fmt.Println("Checking for updates...")
	fmt.Println()

	// Get current version
	currentVersion := Version
	if currentVersion == "" {
		currentVersion = "dev"
	}

	// Fetch latest release from GitHub
	latestRelease, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	latestVersion := strings.TrimPrefix(latestRelease.TagName, "v")
	currentVersionClean := strings.TrimPrefix(currentVersion, "v")

	// Compare versions
	if latestVersion == currentVersionClean {
		fmt.Println("You're running the latest version!")
		fmt.Printf("   Current: %s\n", currentVersion)
		return nil
	}

	// Check if current is newer (dev build)
	if currentVersionClean == "dev" || isNewerVersion(currentVersionClean, latestVersion) {
		fmt.Println("You're running a development or newer version")
		fmt.Printf("   Current: %s\n", currentVersion)
		fmt.Printf("   Latest stable: %s\n", latestVersion)
		return nil
	}

	// New version available
	fmt.Println("A new version is available!")
	fmt.Println()
	fmt.Printf("   Current version: %s\n", currentVersion)
	fmt.Printf("   Latest version:  %s\n", latestVersion)
	fmt.Printf("   Released:        %s\n", latestRelease.PublishedAt.Format("2006-01-02"))
	fmt.Println()

	// Show release notes (first 5 lines)
	if latestRelease.Body != "" {
		fmt.Println("Release Notes:")
		lines := strings.Split(latestRelease.Body, "\n")
		maxLines := 5
		if len(lines) < maxLines {
			maxLines = len(lines)
		}
		for i := 0; i < maxLines; i++ {
			fmt.Printf("   %s\n", lines[i])
		}
		if len(lines) > maxLines {
			fmt.Println("   ...")
		}
		fmt.Println()
	}

	// If check-only mode, stop here
	if updateCheckOnly {
		fmt.Printf("Download: %s\n", latestRelease.HTMLURL)
		return nil
	}

	// Show update instructions based on installation method
	showUpdateInstructions(latestRelease)

	return nil
}

func getLatestRelease() (*GitHubRelease, error) {
	url := "https://api.github.com/repos/Swif7ify/Obelisk-CLI/releases/latest"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Obelisk-CLI")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, err
	}

	return &release, nil
}

func isNewerVersion(current, latest string) bool {
	// Simple version comparison (assumes semantic versioning)
	currentParts := strings.Split(current, ".")
	latestParts := strings.Split(latest, ".")

	for i := 0; i < len(currentParts) && i < len(latestParts); i++ {
		var currentNum, latestNum int
		fmt.Sscanf(currentParts[i], "%d", &currentNum)
		fmt.Sscanf(latestParts[i], "%d", &latestNum)

		if currentNum > latestNum {
			return true
		} else if currentNum < latestNum {
			return false
		}
	}

	return len(currentParts) > len(latestParts)
}

func showUpdateInstructions(release *GitHubRelease) {
	fmt.Println("How to Update:")
	fmt.Println()

	// Detect installation method and show appropriate instructions
	switch runtime.GOOS {
	case "windows":
		fmt.Println("   Option 1: MSI Installer (Recommended)")
		fmt.Printf("   1. Download: %s\n", getAssetURL(release, "msi"))
		fmt.Println("   2. Run the installer")
		fmt.Println("   3. Follow the wizard")
		fmt.Println()
		fmt.Println("   Option 2: Chocolatey")
		fmt.Println("   choco upgrade obelisk-cli")
		fmt.Println()
		fmt.Println("   Option 3: Winget")
		fmt.Println("   winget upgrade OneDev.ObeliskCLI")
		fmt.Println()
		fmt.Println("   Option 4: Manual")
		fmt.Printf("   1. Download: %s\n", getAssetURL(release, "exe"))
		fmt.Println("   2. Replace your current obelisk.exe")

	case "darwin":
		fmt.Println("   Option 1: Homebrew (Recommended)")
		fmt.Println("   brew upgrade obelisk-cli")
		fmt.Println()
		fmt.Println("   Option 2: Manual")
		arch := runtime.GOARCH
		if arch == "amd64" {
			fmt.Printf("   1. Download: %s\n", getAssetURL(release, "darwin-amd64"))
		} else {
			fmt.Printf("   1. Download: %s\n", getAssetURL(release, "darwin-arm64"))
		}
		fmt.Println("   2. Replace your current obelisk binary")

	case "linux":
		fmt.Println("   Option 1: Homebrew")
		fmt.Println("   brew upgrade obelisk-cli")
		fmt.Println()
		fmt.Println("   Option 2: Manual")
		fmt.Printf("   1. Download: %s\n", getAssetURL(release, "linux-amd64"))
		fmt.Println("   2. Replace your current obelisk binary")
	}

	fmt.Println()
	fmt.Printf("Release Page: %s\n", release.HTMLURL)
}

func getAssetURL(release *GitHubRelease, keyword string) string {
	for _, asset := range release.Assets {
		if strings.Contains(strings.ToLower(asset.Name), keyword) {
			return asset.BrowserDownloadURL
		}
	}
	return release.HTMLURL
}

// Made with Bob
