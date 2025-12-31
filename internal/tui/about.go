package tui

import (
	"database/sql"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type AboutModel struct {
	database        *sql.DB
	isOpen          bool
	width           int
	height          int
	viewport        viewport.Model
	glamourRenderer *glamour.TermRenderer
	ready           bool
}

func NewAboutModel(database *sql.DB) AboutModel {
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(72),
	)

	return AboutModel{
		database:        database,
		isOpen:          false,
		glamourRenderer: renderer,
		ready:           false,
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
		switch msg.String() {
		case "q", "esc":
			m.isOpen = false
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m AboutModel) View() string {
	if !m.isOpen {
		return ""
	}

	header := lipgloss.NewStyle().
		Foreground(ColorCyan).
		Bold(true).
		Align(lipgloss.Center).
		Width(80).
		Render("ABOUT ROOTCAMP")

	footer := lipgloss.NewStyle().
		Foreground(ColorGray).
		Italic(true).
		Align(lipgloss.Center).
		Width(80).
		Render("↑/↓ to scroll • q/esc to return to menu")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		m.viewport.View(),
		"",
		footer,
	)

	modal := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(AccentBlue).
		Padding(1, 2).
		Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
		lipgloss.WithWhitespaceChars("░"),
		lipgloss.WithWhitespaceForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}),
	)
}

func (m *AboutModel) Open(width, height int) tea.Cmd {
	m.width = width
	m.height = height
	m.isOpen = true

	aboutContent := `# What is RootCamp?

Most modern interfaces are designed to hide how things actually work. We're here to do the opposite.

**RootCamp** is a hands-on training environment for learning terminal commands. Not through reading or watching—through _doing_. We believe the best way to learn the command line is by actually using it in real, practical scenarios.

## Why should you care?

The terminal isn't just for developers. It's the most powerful interface on your machine. Whether you're managing files, automating tasks, debugging systems, or just trying to understand what's happening under the hood—knowing your way around the command line gives you control.

But here's the problem: most people learn by trial and error on their actual system, which can be scary. What if you delete something important? What if you mess up permissions? What if you get lost in the directory tree?

That's why we built RootCamp.

## The Learning Loop

Every lesson in RootCamp follows a simple pattern:

1. **Learn**: We explain what a command does, why it exists, and when to use it
2. **Practice**: We spin up a real, isolated sandbox environment (` + "`/tmp/rootcamp-{uuid}/`" + `)
3. **Explore**: You use the actual command in a safe space with real files and directories
4. **Validate**: Find the hidden secret code to prove you understand the concept

No multiple choice. No simulated terminal. You're using the _real_ tools, just in a safe playground.

## Safety First

Everything happens in temporary directories under ` + "`/tmp`" + `. When you exit a lesson, the sandbox is automatically cleaned up. No clutter, no residue, no accidentally breaking your system.

You can experiment, make mistakes, and learn without fear.

## Built With

- **Go 1.2x** - Fast, compiled, runs everywhere
- **Bubble Tea** - Modern TUI framework with real-time rendering
- **Lip Gloss** - Beautiful terminal styling
- **Glamour** - Markdown rendering that doesn't suck
- **SQLite** - Track your progress locally
- **Harmonica** - Smooth spring animations

## The Philosophy

We're not trying to replace man pages or cheat sheets. We're building muscle memory.

Reading about ` + "`cd`" + ` is one thing. Actually navigating a nested directory structure to find a hidden file? That's how you learn.

RootCamp is for anyone who wants to feel comfortable in the terminal—whether you're just starting out, or you've been using it for years but want to fill in the gaps.

---

**Version**: 0.1.0
**License**: MIT
**Built by**: Developers who believe the terminal is still the best interface ever created`

	renderedContent, err := m.glamourRenderer.Render(aboutContent)
	if err != nil {
		renderedContent = aboutContent
	}

	vpWidth := 76
	vpHeight := height - 12

	vp := viewport.New(vpWidth, vpHeight)
	vp.SetContent(strings.TrimSpace(renderedContent))

	m.viewport = vp
	m.ready = true

	return nil
}

func (m *AboutModel) Close() {
	m.isOpen = false
}

func (m AboutModel) IsOpen() bool {
	return m.isOpen
}
