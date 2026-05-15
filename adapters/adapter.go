package adapters

import "github.com/Swif7ify/Obelisk-CLI/internal/scanner"

// Adapter defines the interface for framework-specific rules.
type Adapter interface {
	// Name returns the adapter display name.
	Name() string

	// NamingRules returns naming conventions for this framework.
	NamingRules() scanner.NamingRules

	// IgnorePatterns returns additional .gitignore patterns for this framework.
	IgnorePatterns() []string

	// ExpectedDirs returns directories expected in a well-structured project.
	ExpectedDirs() []string

	// Description returns a short description of the framework.
	Description() string
}
