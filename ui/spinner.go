package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// SpinnerPhase represents what the scanner is currently doing.
type SpinnerPhase int

const (
	PhaseSecrets SpinnerPhase = iota
	PhaseGitignore
	PhaseDependencies
	PhaseNaming
	PhaseImports
	PhaseAI
	PhaseDone
)

var phaseLabels = map[SpinnerPhase]string{
	PhaseSecrets:      "Scanning for secrets...",
	PhaseGitignore:    "Validating .gitignore...",
	PhaseDependencies: "Auditing dependencies...",
	PhaseNaming:       "Checking naming conventions...",
	PhaseImports:      "Analyzing imports...",
	PhaseAI:           "Consulting AI brain...",
	PhaseDone:         "Done!",
}

// SpinnerModel is the Bubble Tea model for the scan progress spinner.
type SpinnerModel struct {
	Phase SpinnerPhase
	frame int
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// SpinnerTickMsg triggers a spinner animation frame.
type SpinnerTickMsg struct{}

// NewSpinner creates a new spinner model.
func NewSpinner() SpinnerModel {
	return SpinnerModel{Phase: PhaseSecrets, frame: 0}
}

func (m SpinnerModel) Init() tea.Cmd {
	return nil
}

func (m SpinnerModel) Update(msg tea.Msg) (SpinnerModel, tea.Cmd) {
	switch msg.(type) {
	case SpinnerTickMsg:
		m.frame = (m.frame + 1) % len(spinnerFrames)
	}
	return m, nil
}

func (m SpinnerModel) View() string {
	if m.Phase == PhaseDone {
		return SuccessStyle.Render("✓ Scan complete!")
	}

	frame := spinnerFrames[m.frame]
	label := phaseLabels[m.Phase]

	// Progress bar
	total := int(PhaseDone)
	current := int(m.Phase)
	barWidth := 30
	filled := current * barWidth / total
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	return fmt.Sprintf(
		"\n  %s %s\n  %s %s\n",
		SubtitleStyle.Render(frame),
		label,
		MutedStyle.Render("["+bar+"]"),
		MutedStyle.Render(fmt.Sprintf("%d/%d", current, total)),
	)
}
