package adapters

import "github.com/Swif7ify/Obelisk-CLI/internal/scanner"

// JavaScriptAdapter provides rules for JS/TS/React/Next.js projects.
type JavaScriptAdapter struct{}

func (a *JavaScriptAdapter) Name() string {
	return "JavaScript/TypeScript"
}

func (a *JavaScriptAdapter) NamingRules() scanner.NamingRules {
	return scanner.DefaultJSNamingRules()
}

func (a *JavaScriptAdapter) IgnorePatterns() []string {
	return []string{
		"node_modules/",
		".next/",
		"dist/",
		"build/",
		".env",
		".env.local",
		".env.*.local",
		"coverage/",
		".nyc_output/",
		"*.tsbuildinfo",
		".turbo/",
	}
}

func (a *JavaScriptAdapter) ExpectedDirs() []string {
	return []string{"src", "public"}
}

func (a *JavaScriptAdapter) Description() string {
	return "JavaScript/TypeScript project (React, Next.js, Node.js)"
}
