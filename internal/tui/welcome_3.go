package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

// System checks for boot sequence
var systemChecksV3 = []string{
	"Started Set console font and keymap.",
	"Started Tell Plymouth To Write Out Runtime Data.",
	"Started Create Volatile Files and Directories.",
	"Started Create final runtime dir for shutdown pivot root.",
	"Started Rebuild failed boot detection.",
	"Starting Network Time Synchronization...",
	"Started Authentication service for virtual machines hosted on VMware.",
	"Starting Update UTMP about System Boot/Shutdown...",
	"Started Update UTMP about System Boot/Shutdown.",
	"Started Network Time Synchronization.",
	"Reached target System Time Synchronized.",
	"Started Load AppArmor profiles.",
	"Started Sandbox Environment Initialization.",
	"Started Load Lab Environment Kernel Modules.",
	"Started SQLite Progress Database Service.",
	"Started Lesson Content Provisioning Service.",
	"Starting Initial Sandbox Provisioning...",
	"Reached target RootCamp Training Environment Ready.",
}

// File tree structure that will be revealed
type fileNode struct {
	name     string
	offset   float64
	velocity float64
	spring   harmonica.Spring
	revealed bool
}

const (
	phaseBootSequence = iota
	phaseProvisioning
	phaseComplete
)

type Welcome3Model struct {
	// Boot sequence
	currentCheck int
	bootComplete bool

	// Provisioning phase
	phase       int
	progress    int
	fileNodes   []fileNode
	currentFile int

	// Layout
	width  int
	height int

	// Markdown renderer
	glamourRenderer *glamour.TermRenderer
}

type (
	bootCheckMsg     int
	provisionTickMsg time.Time
	fileRevealMsg    int
)

func NewWelcome3Model() Welcome3Model {
	// Initialize Glamour renderer for markdown
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(40),
	)

	// Initialize file nodes with spring physics
	files := []fileNode{
		{
			name:     "/tmp/rootcamp-x82z/",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "├── bin/",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "│   ├── rootcamp",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "│   └── lab-runner",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "├── etc/",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "│   ├── config.yaml",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "│   └── permissions.conf",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "├── var/",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "│   ├── cache/",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "│   └── run/",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "├── home/",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "│   └── student/",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "│       └── .bashrc",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "├── tmp/",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "│   └── workspace/",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "├── .ghost_dir/",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "├── secrets.txt",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "└── logs/",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "    ├── session.log",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
		{
			name:     "    └── errors.log",
			offset:   -30.0,
			velocity: 0.0,
			spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4),
		},
	}

	return Welcome3Model{
		phase:           phaseBootSequence,
		glamourRenderer: renderer,
		fileNodes:       files,
	}
}

func (m Welcome3Model) Init() tea.Cmd {
	return tea.Batch(
		tickForBootCheck(),
		tickForProvisionAnimation(),
	)
}

func tickForBootCheck() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return bootCheckMsg(0)
	})
}

func tickForProvisionAnimation() tea.Cmd {
	return tea.Tick(16*time.Millisecond, func(t time.Time) tea.Msg {
		return provisionTickMsg(t)
	})
}

func tickForFileReveal() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return fileRevealMsg(0)
	})
}

func (m Welcome3Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case bootCheckMsg:
		if m.phase == phaseBootSequence {
			if m.currentCheck < len(systemChecksV3) {
				m.currentCheck++
				if m.currentCheck >= len(systemChecksV3) {
					m.bootComplete = true
					// Wait a moment then transition to provisioning
					return m, tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
						return provisionTickMsg(t)
					})
				}
				return m, tickForBootCheck()
			}
		}

	case provisionTickMsg:
		if m.bootComplete && m.phase == phaseBootSequence {
			// Transition to provisioning phase
			m.phase = phaseProvisioning
			return m, tea.Batch(
				tickForFileReveal(),
				tickForProvisionAnimation(),
			)
		}

		if m.phase == phaseProvisioning {
			// Update progress
			if m.progress < 100 {
				m.progress += 1
				if m.progress > 100 {
					m.progress = 100
					m.phase = phaseComplete
				}
			}

			return m, tickForProvisionAnimation()
		}

	case fileRevealMsg:
		if m.phase == phaseProvisioning && m.currentFile < len(m.fileNodes) {
			m.fileNodes[m.currentFile].revealed = true
			m.currentFile++
			if m.currentFile < len(m.fileNodes) {
				return m, tickForFileReveal()
			}
		}
	}

	return m, nil
}

func (m Welcome3Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	if m.phase == phaseBootSequence {
		return m.renderBootSequence()
	}

	return m.renderProvisioningView()
}

func (m Welcome3Model) renderBootSequence() string {
	var output strings.Builder

	// Green color for OK status (like systemd)
	okStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true)

	// White/gray for message text
	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	// Dim style for "Starting" messages
	startingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#999999"))

	output.WriteString("\n")

	for i := 0; i < m.currentCheck && i < len(systemChecksV3); i++ {
		message := systemChecksV3[i]

		// Check if this is a "Starting" message (show without OK status yet)
		if strings.HasPrefix(message, "Starting") {
			line := fmt.Sprintf("         %s\n", startingStyle.Render(message))
			output.WriteString(line)
		} else {
			// Show with green OK status
			line := fmt.Sprintf("  %s %s\n",
				okStyle.Render("[ OK ]"),
				messageStyle.Render(message))
			output.WriteString(line)
		}
	}

	return output.String()
}

func (m Welcome3Model) renderProvisioningView() string {
	// Three columns: File tree, Menu, Architect's Log
	fileTree := m.renderFileTree()
	menu := m.renderMenu()
	architectLog := m.renderArchitectLog()

	// Progress bar
	progressBar := m.renderProgressBar()

	// Create the 3-column layout
	leftWidth := 50
	rightWidth := 50
	middleWidth := m.width - leftWidth - rightWidth - 10 // Fill remaining width, account for borders

	panelHeight := m.height - 12 // Leave room for header, footer, progress

	// Style the left panel (Sandbox Structure)
	leftStyle := lipgloss.NewStyle().
		Width(leftWidth).
		Height(panelHeight).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5FB3FF"))

	// Style the middle panel (Menu)
	middleStyle := lipgloss.NewStyle().
		Width(middleWidth).
		Height(panelHeight).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FFB86C"))

	// Style the right panel (Architect's Log)
	rightStyle := lipgloss.NewStyle().
		Width(rightWidth).
		Height(panelHeight).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#9D7CFF"))

	left := leftStyle.Render(fileTree)
	middle := middleStyle.Render(menu)
	right := rightStyle.Render(architectLog)

	// Join all three columns horizontally
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, left, middle, right)

	// Create header
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7DCFFF")).
		Bold(true).
		Width(m.width).
		Align(lipgloss.Center)

	header := headerStyle.Render("ROOT CAMP v0.1")

	// Create footer with progress
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Italic(true)

	footer := footerStyle.Render("(q) to exit")

	// Center progress bar and footer
	centeredProgressBar := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(progressBar)

	centeredFooter := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(footer)

	// Assemble everything
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		header,
		"",
		mainContent,
		"",
		centeredProgressBar,
		"",
		centeredFooter,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m Welcome3Model) renderMenu() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1A1B26")).
		Background(lipgloss.Color("#FFB86C")).
		Bold(true).
		Padding(0, 1)

	optionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7DCFFF"))

	disabledStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Italic(true)

	// Stubbed menu options
	menuOptions := []struct {
		key     string
		label   string
		enabled bool
	}{
		{"1", "Start Training", false},
		{"2", "View Lessons", false},
		{"3", "Lab Environment", false},
		{"4", "Progress Tracker", false},
		{"5", "Settings", false},
		{"q", "Exit", true},
	}

	// Calculate available width for the menu panel
	leftWidth := 50
	rightWidth := 50
	middleWidth := m.width - leftWidth - rightWidth - 10

	// Center the title separately
	title := titleStyle.Render("MAIN MENU")
	centeredTitle := lipgloss.Place(
		middleWidth-4, // Account for padding
		1,
		lipgloss.Center,
		lipgloss.Center,
		title,
	)

	// Build menu items with left alignment
	var menuItems []string
	for _, opt := range menuOptions {
		var line string
		if opt.enabled {
			line = optionStyle.Render(fmt.Sprintf("[%s] %s", opt.key, opt.label))
		} else {
			line = disabledStyle.Render(fmt.Sprintf("[%s] %s", opt.key, opt.label))
		}
		menuItems = append(menuItems, line)
	}

	// Join menu items with LEFT alignment
	leftAlignedItems := lipgloss.JoinVertical(lipgloss.Left, menuItems...)

	// Center the left-aligned items block
	centeredItems := lipgloss.Place(
		middleWidth-4,
		len(menuItems),
		lipgloss.Center,
		lipgloss.Top,
		leftAlignedItems,
	)

	// Combine centered title with centered (but internally left-aligned) menu items
	return lipgloss.JoinVertical(
		lipgloss.Left,
		centeredTitle,
		"",
		centeredItems,
	)
}

func (m Welcome3Model) renderFileTree() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1A1B26")).
		Background(lipgloss.Color("#5FB3FF")).
		Bold(true).
		Padding(0, 1)

	fileStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9ECE6A"))

	// Center the title
	title := titleStyle.Render("SANDBOX STRUCTURE")
	leftWidth := 50
	centeredTitle := lipgloss.Place(
		leftWidth-4, // Account for padding
		1,
		lipgloss.Center,
		lipgloss.Center,
		title,
	)

	// Build file tree
	var fileList strings.Builder
	for _, node := range m.fileNodes {
		if node.revealed {
			fileList.WriteString(fileStyle.Render(node.name))
			fileList.WriteString("\n")
		}
	}

	return centeredTitle + "\n\n" + fileList.String()
}

func (m Welcome3Model) renderArchitectLog() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1A1B26")).
		Background(lipgloss.Color("#9D7CFF")).
		Bold(true).
		Padding(0, 1)

	// Center the title
	title := titleStyle.Render("THE ARCHITECT'S LOG")
	rightWidth := 50
	centeredTitle := lipgloss.Place(
		rightWidth-4, // Account for padding
		1,
		lipgloss.Center,
		lipgloss.Center,
		title,
	)

	markdown := `
In **1965**, the Multics operating system introduced the revolutionary concept of a hierarchical directory structure.

Before this breakthrough, data storage was a _flat pile of magnetic tape_—no organization, no hierarchy, just sequential blocks.

The directory changed everything. It gave us:
- **Namespaces** for file organization
- **Paths** to navigate data
- **Permissions** to control access

Today, every terminal command you run traverses this tree. You're not just using the filesystem—you're walking through history.`

	// Render with Glamour
	rendered, err := m.glamourRenderer.Render(markdown)
	if err != nil {
		return centeredTitle + "\n\n" + markdown // Fallback to plain text
	}

	return centeredTitle + "\n\n" + strings.TrimSpace(rendered)
}

func (m Welcome3Model) renderProgressBar() string {
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Italic(true)

	barWidth := 50
	filled := int(float64(barWidth) * float64(m.progress) / 100.0)
	empty := barWidth - filled

	filledStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9ECE6A")).
		Bold(true)

	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#414868"))

	bar := filledStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", empty))

	label := labelStyle.Render(fmt.Sprintf("Status: [Provisioning Sandbox...] %d%%", m.progress))

	return lipgloss.JoinVertical(
		lipgloss.Center,
		label,
		bar,
	)
}
