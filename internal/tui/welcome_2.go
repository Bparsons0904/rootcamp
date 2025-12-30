package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Welcome2Model struct {
	stopwatch stopwatch.Model
	width     int
	height    int
}

func NewWelcome2Model() Welcome2Model {
	return Welcome2Model{
		stopwatch: stopwatch.New(),
	}
}

func (m Welcome2Model) Init() tea.Cmd {
	return m.stopwatch.Init()
}

func (m Welcome2Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

func (m Welcome2Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Generate the ASCII grid background (dim directory tree)
	background := m.generateDirectoryGrid()

	// Create the vibrant title with Tokyo Night gradient
	title := m.createTitle()

	// Create the stopwatch display
	stopwatchDisplay := m.createStopwatchDisplay()

	// Layer the components using Lip Gloss positioning
	// Center everything
	centeredTitle := lipgloss.Place(
		m.width,
		m.height-5,
		lipgloss.Center,
		lipgloss.Center,
		title,
	)

	// Position stopwatch at bottom
	centeredStopwatch := lipgloss.Place(
		m.width,
		3,
		lipgloss.Center,
		lipgloss.Center,
		stopwatchDisplay,
	)

	// Dim the background
	dimBackground := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1a1b26")).
		Render(background)

	// Layer background and foreground
	output := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Left,
		lipgloss.Top,
		dimBackground,
	)

	// Overlay the centered title on top of the background
	var result strings.Builder
	backgroundLines := strings.Split(output, "\n")
	titleLines := strings.Split(centeredTitle, "\n")

	maxLines := len(backgroundLines)
	if len(titleLines) > maxLines {
		maxLines = len(titleLines)
	}

	// Merge background and title (title overwrites background where it has content)
	for i := 0; i < maxLines; i++ {
		var bgLine, titleLine string
		if i < len(backgroundLines) {
			bgLine = backgroundLines[i]
		}
		if i < len(titleLines) {
			titleLine = titleLines[i]
		}

		// If title line has content, use it; otherwise use background
		if strings.TrimSpace(titleLine) != "" {
			result.WriteString(titleLine)
		} else {
			result.WriteString(bgLine)
		}
		result.WriteString("\n")
	}

	// Add stopwatch at the bottom
	result.WriteString(centeredStopwatch)

	return result.String()
}

func (m Welcome2Model) generateDirectoryGrid() string {
	var grid strings.Builder

	// Create a directory tree structure
	tree := []string{
		"rootcamp/",
		"├── bin/",
		"│   ├── rootcamp",
		"│   └── lab-runner",
		"├── lessons/",
		"│   ├── basics/",
		"│   │   ├── cd/",
		"│   │   ├── ls/",
		"│   │   ├── pwd/",
		"│   │   └── mkdir/",
		"│   ├── intermediate/",
		"│   │   ├── grep/",
		"│   │   ├── find/",
		"│   │   ├── chmod/",
		"│   │   └── chown/",
		"│   └── advanced/",
		"│       ├── awk/",
		"│       ├── sed/",
		"│       └── xargs/",
		"├── labs/",
		"│   ├── sandbox/",
		"│   └── tmp/",
		"├── db/",
		"│   └── progress.db",
		"└── kernel/",
		"    ├── init",
		"    ├── scheduler",
		"    └── syscalls",
	}

	// Repeat the tree to fill the screen vertically
	linesNeeded := m.height
	currentLine := 0

	for currentLine < linesNeeded {
		for _, line := range tree {
			if currentLine >= linesNeeded {
				break
			}
			grid.WriteString(line)

			// Pad the rest of the line with spaces to fill width
			remaining := m.width - len(line)
			if remaining > 0 {
				grid.WriteString(strings.Repeat(" ", remaining))
			}
			grid.WriteString("\n")
			currentLine++
		}
	}

	return grid.String()
}

func (m Welcome2Model) createTitle() string {
	// Create gradient title using Tokyo Night colors
	titleParts := []struct {
		text  string
		color string
	}{
		{"█▀█ ", "#F7768E"}, // Red
		{"█▀█ ", "#FF9E64"}, // Orange
		{"█▀█ ", "#E0AF68"}, // Yellow
		{"▀█▀ ", "#9ECE6A"}, // Green
		{"█▀▀ ", "#7DCFFF"}, // Cyan
		{"▄▀█ ", "#7AA2F7"}, // Blue
		{"█▀▄▀█ ", "#BB9AF7"}, // Purple
		{"█▀█", "#C0CAF5"},   // Light purple
	}

	var titleLine strings.Builder
	for _, part := range titleParts {
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color(part.color)).
			Bold(true)
		titleLine.WriteString(style.Render(part.text))
	}

	titleText := titleLine.String()

	subtitle1Style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7DCFFF")).
		Bold(true)
	subtitle1 := subtitle1Style.Render("ROOT CAMP")

	subtitle2Style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Italic(true)
	subtitle2 := subtitle2Style.Render("Terminal Mastery • System Control • Production Ready")

	taglineStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9ECE6A")).
		Bold(true).
		Italic(true)
	tagline := taglineStyle.Render("Master the command line. Control the kernel.")

	// Build the complete title box
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		titleText,
		"",
		subtitle1,
		"",
		subtitle2,
		"",
		"",
		tagline,
		"",
	)

	// Create a box with border
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7AA2F7")).
		Padding(1, 3).
		Align(lipgloss.Center)

	return boxStyle.Render(content)
}

func (m Welcome2Model) createStopwatchDisplay() string {
	timerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E0AF68")).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Italic(true)

	reminderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F7768E")).
		Italic(true)

	timer := timerStyle.Render(m.stopwatch.View())
	label := labelStyle.Render("Session Time")
	reminder := reminderStyle.Render("In production, time is the one resource you can't chmod")

	display := lipgloss.JoinVertical(
		lipgloss.Center,
		label,
		timer,
		"",
		reminder,
		"",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#565F89")).Italic(true).Render("Press 'q' to exit"),
	)

	return display
}
