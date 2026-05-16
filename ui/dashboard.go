package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Swif7ify/Obelisk-CLI/internal/ai"
	"github.com/Swif7ify/Obelisk-CLI/internal/scanner"
)

// DashboardState represents the current state of the dashboard.
type DashboardState int

const (
	StateScanning DashboardState = iota
	StateResults
)

// ScanCompleteMsg is sent when scanning is complete.
type ScanCompleteMsg struct {
	Result *scanner.ScanResult
	Report *ai.HealthReport
	Err    error
}

// DashboardModel is the main Bubble Tea model for the health check dashboard.
type DashboardModel struct {
	State   DashboardState
	Spinner SpinnerModel
	Result  *scanner.ScanResult
	Report  *ai.HealthReport
	Err     error
	Width   int
	Height  int
}

// NewDashboard creates a new dashboard model.
func NewDashboard() DashboardModel {
	return DashboardModel{
		State:   StateScanning,
		Spinner: NewSpinner(),
		Width:   80,
		Height:  40,
	}
}

func (m DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
	)
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case SpinnerTickMsg:
		if m.State == StateScanning {
			m.Spinner, _ = m.Spinner.Update(msg)
			return m, tickCmd()
		}

	case ScanCompleteMsg:
		m.State = StateResults
		m.Result = msg.Result
		m.Report = msg.Report
		m.Err = msg.Err
		m.Spinner.Phase = PhaseDone
	}

	return m, nil
}

func (m DashboardModel) View() string {
	var s string

	s += RenderBanner()

	switch m.State {
	case StateScanning:
		s += m.Spinner.View()

	case StateResults:
		if m.Err != nil {
			s += ErrorStyle.Render("✗ Error: " + m.Err.Error()) + "\n"
		} else if m.Report != nil && m.Result != nil {
			s += RenderScoreCard(m.Report) + "\n"
		}
	}

	s += "\n" + MutedStyle.Render("Press 'q' to quit") + "\n"
	return s
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return SpinnerTickMsg{}
	})
}
