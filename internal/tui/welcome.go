package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

// System checks that will scroll during boot sequence
var systemChecks = []struct {
	name   string
	status string
}{
	{"FS_SANDBOX_READY", "OK"},
	{"SQLITE_DB_LINKED", "OK"},
	{"LESSON_MODULES_LOADED", "OK"},
	{"TTY_INTERFACE_INIT", "OK"},
	{"USER_PROGRESS_TRACKER", "OK"},
	{"LAB_ENVIRONMENT_CHECK", "OK"},
	{"COMMAND_PARSER_READY", "OK"},
	{"KERNEL_PRIVILEGES_SET", "OK"},
}

type WelcomeModel struct {
	currentCheck int
	checksComplete bool
	spring harmonica.Spring
	boxOffset float64
	boxVelocity float64
	width int
	height int
}

type checkCompleteMsg int
type animationTickMsg time.Time

func NewWelcomeModel() WelcomeModel {
	return WelcomeModel{
		currentCheck: 0,
		checksComplete: false,
		spring: harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.5),
		boxOffset: -50.0,
		boxVelocity: 0.0,
	}
}

func (m WelcomeModel) Init() tea.Cmd {
	return tea.Batch(
		tickForCheck(),
		tickForAnimation(),
	)
}

func tickForCheck() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return checkCompleteMsg(0)
	})
}

func tickForAnimation() tea.Cmd {
	return tea.Tick(16*time.Millisecond, func(t time.Time) tea.Msg {
		return animationTickMsg(t)
	})
}

func (m WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case checkCompleteMsg:
		if m.currentCheck < len(systemChecks) {
			m.currentCheck++
			if m.currentCheck >= len(systemChecks) {
				m.checksComplete = true
			}
			return m, tickForCheck()
		}

	case animationTickMsg:
		if m.checksComplete {
			m.boxOffset, m.boxVelocity = m.spring.Update(m.boxOffset, m.boxVelocity, 0.0)
		}
		return m, tickForAnimation()
	}

	return m, nil
}

func (m WelcomeModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var output strings.Builder

	// Define vibrant, modern color scheme
	checkStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#61FFCA")).
		Bold(true)

	statusOKStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7AA2F7")).
		Bold(true)

	// Render system checks
	checksOutput := strings.Builder{}
	for i := 0; i < m.currentCheck && i < len(systemChecks); i++ {
		check := systemChecks[i]
		line := fmt.Sprintf("  %s... %s\n",
			checkStyle.Render(check.name),
			statusOKStyle.Render(check.status))
		checksOutput.WriteString(line)
	}

	// Calculate vertical position for checks
	checksBlock := checksOutput.String()
	checksHeight := strings.Count(checksBlock, "\n")
	checksPadding := (m.height - checksHeight - 20) / 2
	if checksPadding < 0 {
		checksPadding = 0
	}

	output.WriteString(strings.Repeat("\n", checksPadding))
	output.WriteString(checksBlock)

	// If all checks are complete, show the main box with spring animation
	if m.checksComplete {
		// Create the title with gradient effect (Tokyo Night inspired)
		titleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BB9AF7")).
			Bold(true).
			Italic(true)

		title := titleStyle.Render("R O O T C A M P")

		subtitle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7DCFFF")).
			Render("Terminal Mastery Bootcamp")

		tagline := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9ECE6A")).
			Italic(true).
			Render("From noob to root")

		hookStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F7768E")).
			Bold(true).
			Italic(true)

		hook := hookStyle.Render("The kernel is ready. Are you?")

		// Build the box content
		boxContent := lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			title,
			"",
			subtitle,
			"",
			tagline,
			"",
			"",
			hook,
			"",
		)

		// Create the double-bordered box
		boxStyle := lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#F7768E")).
			Padding(2, 4).
			Align(lipgloss.Center)

		box := boxStyle.Render(boxContent)

		// Apply spring animation offset
		boxLines := strings.Split(box, "\n")
		animatedBox := strings.Builder{}

		// Calculate the offset in lines (convert float offset to int)
		lineOffset := int(m.boxOffset)
		if lineOffset < 0 {
			lineOffset = 0
		}

		// Add vertical offset
		animatedBox.WriteString(strings.Repeat("\n", lineOffset))

		// Center the box horizontally
		for _, line := range boxLines {
			padding := (m.width - lipgloss.Width(line)) / 2
			if padding < 0 {
				padding = 0
			}
			animatedBox.WriteString(strings.Repeat(" ", padding))
			animatedBox.WriteString(line)
			animatedBox.WriteString("\n")
		}

		output.WriteString("\n")
		output.WriteString(animatedBox.String())

		// Add press any key hint at bottom
		hintStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565F89")).
			Italic(true)

		hint := hintStyle.Render("Press 'q' to exit")
		hintPadding := (m.width - lipgloss.Width(hint)) / 2
		if hintPadding < 0 {
			hintPadding = 0
		}

		output.WriteString("\n\n")
		output.WriteString(strings.Repeat(" ", hintPadding))
		output.WriteString(hint)
	}

	return output.String()
}
