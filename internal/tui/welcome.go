package tui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

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

var architectFacts = []string{
	`In **1965**, Multics introduced hierarchical directories.

Before this, data was a _flat pile of magnetic tape_—no hierarchy.

The directory changed everything:
- **Namespaces** for organization
- **Paths** to navigate data
- **Permissions** to control access

Every command you run traverses this tree.`,

	`The **pipe** (|) was invented by Doug McIlroy in **1973**.

His vision: _"Write programs that do one thing well and work together."_

This single character changed software:
- **Composability** of tools
- The Unix **philosophy**
- Simple power: **cat log | grep ERROR**`,

	`In **1991**, Linus Torvalds posted to comp.os.minix:

_"I'm doing a (free) OS (just a hobby)..."_

That hobby became **Linux**:
- Powers **96.3%** of top servers
- Every **Android** device
- **100%** of top 500 supercomputers

A student project, now foundation of the internet.`,

	`**/bin** and **/usr/bin** split: **1971**.

Unix ran out of disk space. The PDP-11 had a **1.5MB** drive. Dennis Ritchie added a second disk at /usr.

Today:
- **/bin** - System binaries
- **/usr/bin** - User programs

Modern Linux carries the ghost of a 1970s disk shortage.`,

	`**chmod** uses octal due to **1974** hardware limits.

Permissions needed **9 bits** (rwxrwxrwx). Octal aligned perfectly:
- **755** = rwxr-xr-x
- **644** = rw-r--r--

Memory was expensive, every bit counted.

We still use octal because _that's how it's always been_.`,

	`The **root** user was never meant to be permanent.

Ken Thompson created it for testing. It was **temporary**.

Instead, it became **immortal**:
- Every Unix system has root
- 50+ years later, still here

The ultimate _"temporary solution"_.`,

	`**Hidden files** (.bashrc) were an accident.

Early **ls** sorted alphabetically. Files starting with **.** sorted first.

Later, someone hid dotfiles to reduce clutter.

Result:
- Configs became "special"
- Pattern became **convention**

Your home is littered with dotfiles from a sorting hack.`,

	`**/etc** means **"et cetera"**—_"and other things."_

Early Unix had:
- **/bin** for binaries
- **/dev** for devices
- **/lib** for libraries

Everything else? **Et cetera.**

The "misc folder" became the backbone of system admin.`,

	`**tty** = **teletypewriter** (1920s hardware).

Early terminals were literal typewriters:
- No screen, just paper
- Type a command, it prints

Modern terminals are **emulators** of 100-year-old machines.

A simulation of a simulation.`,

	`The **$** prompt has military origins.

In the **1960s**, computing cost money per CPU second. **$** reminded users:
_"This costs money"_

Root used **#** (override costs).

Today, free computing, but **$** and **#** remain.`,
}

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
	currentCheck    int
	bootComplete    bool
	phase           int
	progress        int
	fileNodes       []fileNode
	currentFile     int
	width           int
	height          int
	glamourRenderer *glamour.TermRenderer
	selectedFact    string
}

type (
	bootCheckMsg     int
	provisionTickMsg time.Time
	fileRevealMsg    int
)

func NewWelcome3Model() Welcome3Model {
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(40),
	)

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

	randomIndex := rand.Intn(len(architectFacts))

	return Welcome3Model{
		phase:           phaseBootSequence,
		glamourRenderer: renderer,
		fileNodes:       files,
		selectedFact:    architectFacts[randomIndex],
		width:           120, // Default width to prevent "Loading..." hang
		height:          40,  // Default height, will be updated by WindowSizeMsg
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
	output.WriteString("\n")

	for i := 0; i < m.currentCheck && i < len(systemChecksV3); i++ {
		message := systemChecksV3[i]

		if strings.HasPrefix(message, "Starting") {
			line := fmt.Sprintf("         %s\n", BootStartingStyle().Render(message))
			output.WriteString(line)
		} else {
			line := fmt.Sprintf("  %s %s\n",
				BootOKStyle().Render("[ OK ]"),
				BootMessageStyle().Render(message))
			output.WriteString(line)
		}
	}

	return output.String()
}

func (m Welcome3Model) renderProvisioningView() string {
	fileTree := m.renderFileTree()
	menu := m.renderMenu()
	architectLog := m.renderArchitectLog()
	progressBar := m.renderProgressBar()

	leftWidth := 50
	rightWidth := 50
	middleWidth := m.width - leftWidth - rightWidth - 10
	panelHeight := m.height - 10

	left := PanelStyle(leftWidth, panelHeight, ColorBlue).Render(fileTree)
	middle := PanelStyle(middleWidth, panelHeight, ColorOrange).Render(menu)
	right := PanelStyle(rightWidth, panelHeight, ColorPurple).Render(architectLog)

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, left, middle, right)

	header := HeaderStyle(m.width).Render("ROOT CAMP v0.1")
	footer := FooterStyle().Render("(q) to exit")

	centeredProgressBar := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(progressBar)

	centeredFooter := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(footer)

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

	leftWidth := 50
	rightWidth := 50
	middleWidth := m.width - leftWidth - rightWidth - 10

	title := PanelTitleStyle(ColorOrange).Render("MAIN MENU")
	centeredTitle := lipgloss.Place(
		middleWidth-4,
		1,
		lipgloss.Center,
		lipgloss.Center,
		title,
	)

	var menuItems []string
	for _, opt := range menuOptions {
		var line string
		if opt.enabled {
			line = MenuOptionStyle().Render(fmt.Sprintf("[%s] %s", opt.key, opt.label))
		} else {
			line = DisabledOptionStyle().Render(fmt.Sprintf("[%s] %s", opt.key, opt.label))
		}
		menuItems = append(menuItems, line)
	}

	leftAlignedItems := lipgloss.JoinVertical(lipgloss.Left, menuItems...)

	centeredItems := lipgloss.Place(
		middleWidth-4,
		len(menuItems),
		lipgloss.Center,
		lipgloss.Top,
		leftAlignedItems,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		centeredTitle,
		"",
		centeredItems,
	)
}

func (m Welcome3Model) renderFileTree() string {
	leftWidth := 50
	title := PanelTitleStyle(ColorBlue).Render("SANDBOX STRUCTURE")
	centeredTitle := lipgloss.Place(
		leftWidth-4,
		1,
		lipgloss.Center,
		lipgloss.Center,
		title,
	)

	var fileList strings.Builder
	for _, node := range m.fileNodes {
		if node.revealed {
			fileList.WriteString(FileTreeStyle().Render(node.name))
			fileList.WriteString("\n")
		}
	}

	return centeredTitle + "\n\n" + fileList.String()
}

func (m Welcome3Model) renderArchitectLog() string {
	rightWidth := 50
	title := PanelTitleStyle(ColorPurple).Render("THE ARCHITECT'S LOG")
	centeredTitle := lipgloss.Place(
		rightWidth-4,
		1,
		lipgloss.Center,
		lipgloss.Center,
		title,
	)

	rendered, err := m.glamourRenderer.Render(m.selectedFact)
	if err != nil {
		return centeredTitle + "\n\n" + m.selectedFact
	}

	return centeredTitle + "\n\n" + strings.TrimSpace(rendered)
}

func (m Welcome3Model) renderProgressBar() string {
	barWidth := 50
	filled := int(float64(barWidth) * float64(m.progress) / 100.0)
	empty := barWidth - filled

	bar := ProgressBarFilledStyle().Render(strings.Repeat("█", filled)) +
		ProgressBarEmptyStyle().Render(strings.Repeat("░", empty))

	label := ProgressLabelStyle().Render(fmt.Sprintf("Status: [Provisioning Sandbox...] %d%%", m.progress))

	return lipgloss.JoinVertical(
		lipgloss.Center,
		label,
		bar,
	)
}
