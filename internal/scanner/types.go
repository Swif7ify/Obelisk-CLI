package scanner

import "fmt"

// Severity represents the severity level of a finding.
type Severity int

const (
	SeverityInfo Severity = iota
	SeverityWarning
	SeverityError
	SeverityCritical
)

// String returns the human-readable severity label.
func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "INFO"
	case SeverityWarning:
		return "WARNING"
	case SeverityError:
		return "ERROR"
	case SeverityCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// Category represents the category of a finding.
type Category string

const (
	CategorySecurity     Category = "Security"
	CategoryArchitecture Category = "Architecture"
	CategoryQuality      Category = "Quality"
	CategoryDependency   Category = "Dependency"
	CategoryNaming       Category = "Naming"
)

// Finding represents a single issue discovered during scanning.
type Finding struct {
	Category    Category `json:"category"`
	Severity    Severity `json:"severity"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	File        string   `json:"file,omitempty"`
	Line        int      `json:"line,omitempty"`
	Suggestion  string   `json:"suggestion,omitempty"`
}

// String returns a formatted string for the finding.
func (f Finding) String() string {
	loc := ""
	if f.File != "" {
		loc = fmt.Sprintf(" [%s", f.File)
		if f.Line > 0 {
			loc += fmt.Sprintf(":%d", f.Line)
		}
		loc += "]"
	}
	return fmt.Sprintf("[%s] %s: %s%s", f.Severity, f.Category, f.Title, loc)
}

// ScanResult holds all findings from a complete scan.
type ScanResult struct {
	ProjectPath string    `json:"project_path"`
	ProjectType string    `json:"project_type"`
	Findings    []Finding `json:"findings"`
	FileCount   int       `json:"file_count"`
	DirCount    int       `json:"dir_count"`
	DirTree     string    `json:"dir_tree,omitempty"`
}

// CountBySeverity returns the number of findings for a given severity.
func (r *ScanResult) CountBySeverity(s Severity) int {
	count := 0
	for _, f := range r.Findings {
		if f.Severity == s {
			count++
		}
	}
	return count
}

// CountByCategory returns the number of findings for a given category.
func (r *ScanResult) CountByCategory(c Category) int {
	count := 0
	for _, f := range r.Findings {
		if f.Category == c {
			count++
		}
	}
	return count
}

// HasCritical returns true if any critical findings exist.
func (r *ScanResult) HasCritical() bool {
	return r.CountBySeverity(SeverityCritical) > 0
}
