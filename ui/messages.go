package ui

import tea "github.com/charmbracelet/bubbletea"

// PhaseUpdateMsg is sent when the scan engine moves to a new phase.
type PhaseUpdateMsg struct {
	Phase string
}

// ListenForPhaseUpdates creates a command that listens on a channel for phase updates.
func ListenForPhaseUpdates(c chan string) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-c
		if !ok {
			return nil // Channel closed
		}
		return PhaseUpdateMsg{Phase: msg}
	}
}

// MapPhaseStringToSpinnerPhase converts the string phase from engine.Run into a SpinnerPhase enum.
func MapPhaseStringToSpinnerPhase(phase string) SpinnerPhase {
	switch phase {
	case "Detecting project type...", "Scanning for secrets...":
		return PhaseSecrets
	case "Validating .gitignore...":
		return PhaseGitignore
	case "Auditing dependencies...":
		return PhaseDependencies
	case "Checking naming conventions...":
		return PhaseNaming
	case "Analyzing imports...", "Analyzing cyclomatic complexity...", "Tracking technical debt...", "Running linters...":
		return PhaseImports
	case "Consulting AI brain...":
		return PhaseAI
	default:
		return PhaseSecrets
	}
}
