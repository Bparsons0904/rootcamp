package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	appTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF6B9D")).
			Background(lipgloss.Color("#1a1a1a")).
			Padding(0, 1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			MarginBottom(1)

	sectionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#F25D94")).
				MarginTop(1).
				MarginBottom(1)

	textStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDDDDD"))

	codeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7EC8E3")).
			Background(lipgloss.Color("#2D2D2D")).
			Padding(0, 1)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true).
				PaddingLeft(2)

	unselectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				PaddingLeft(4)

	completedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	incompleteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#2D2D2D")).
			Padding(0, 1).
			MarginTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			MarginTop(1)

	taskStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Background(lipgloss.Color("#2D2D2D")).
			Padding(1).
			MarginTop(1).
			MarginBottom(1)

	pathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7EC8E3")).
			Bold(true)
)
