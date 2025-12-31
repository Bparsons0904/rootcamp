package tui

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AboutModel struct {
	database *sql.DB
	isOpen   bool
	width    int
	height   int
}

func NewAboutModel(database *sql.DB) AboutModel {
	return AboutModel{
		database: database,
		isOpen:   false,
	}
}

func (m AboutModel) Init() tea.Cmd {
	return nil
}

func (m *AboutModel) Update(msg tea.Msg) (*AboutModel, tea.Cmd) {
	if !m.isOpen {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" || msg.String() == "q" {
			m.isOpen = false
			return m, nil
		}
	}

	return m, nil
}

func (m AboutModel) View() string {
	if !m.isOpen {
		return ""
	}

	content := lipgloss.NewStyle().
		Width(80).
		Height(20).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(AccentBlue).
		Padding(2, 4).
		Render("About Root Camp\n\n(Coming Soon)\n\nPress 'q' or 'esc' to return to menu")

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
		lipgloss.WithWhitespaceChars("░"),
		lipgloss.WithWhitespaceForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}),
	)
}

func (m *AboutModel) Open(width, height int) tea.Cmd {
	m.width = width
	m.height = height
	m.isOpen = true
	return nil
}

func (m *AboutModel) Close() {
	m.isOpen = false
}

func (m AboutModel) IsOpen() bool {
	return m.isOpen
}
