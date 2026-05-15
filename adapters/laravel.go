package adapters

import (
	"regexp"

	"github.com/Swif7ify/Obelisk-CLI/internal/scanner"
)

// LaravelAdapter provides stub rules for Laravel/PHP projects (future expansion).
type LaravelAdapter struct{}

func (a *LaravelAdapter) Name() string {
	return "Laravel (PHP)"
}

func (a *LaravelAdapter) NamingRules() scanner.NamingRules {
	return scanner.NamingRules{
		ComponentPattern: regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*\.php$`),
		AssetPattern:     regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*\.[a-z]+$`),
		ComponentDirs:    []string{"app", "Http", "Controllers", "Models"},
		AssetDirs:        []string{"resources", "public"},
	}
}

func (a *LaravelAdapter) IgnorePatterns() []string {
	return []string{"vendor/", ".env", "storage/", "bootstrap/cache/"}
}

func (a *LaravelAdapter) ExpectedDirs() []string {
	return []string{"app", "resources", "routes", "database", "public"}
}

func (a *LaravelAdapter) Description() string {
	return "Laravel PHP project (support coming soon)"
}
