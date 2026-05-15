package adapters

import (
	"regexp"

	"github.com/Swif7ify/Obelisk-CLI/internal/scanner"
)

// GolangAdapter provides stub rules for Go projects (future expansion).
type GolangAdapter struct{}

func (a *GolangAdapter) Name() string {
	return "Go (Golang)"
}

func (a *GolangAdapter) NamingRules() scanner.NamingRules {
	return scanner.NamingRules{
		ComponentPattern: regexp.MustCompile(`^[a-z][a-z0-9_]*\.go$`),
		AssetPattern:     regexp.MustCompile(`^[a-z][a-z0-9_]*\.[a-z]+$`),
		ComponentDirs:    []string{"cmd", "internal", "pkg"},
		AssetDirs:        []string{"assets", "static"},
	}
}

func (a *GolangAdapter) IgnorePatterns() []string {
	return []string{"bin/", "vendor/", ".env", "*.exe"}
}

func (a *GolangAdapter) ExpectedDirs() []string {
	return []string{"cmd", "internal"}
}

func (a *GolangAdapter) Description() string {
	return "Go project (support coming soon)"
}
