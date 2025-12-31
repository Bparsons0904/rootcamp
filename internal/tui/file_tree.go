package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

type fileNode struct {
	name     string
	offset   float64
	velocity float64
	spring   harmonica.Spring
	revealed bool
}

type fileRevealMsg int

type FileTreeModel struct {
	nodes       []fileNode
	currentFile int
	complete    bool
	width       int
}

func NewFileTreeModel(width int, skipAnimations bool) FileTreeModel {
	files := []fileNode{
		{name: "/tmp/rootcamp-x82z/", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "├── bin/", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "│   ├── rootcamp", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "│   └── lab-runner", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "├── etc/", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "│   ├── config.yaml", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "│   └── permissions.conf", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "├── var/", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "│   ├── cache/", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "│   └── run/", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "├── home/", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "│   └── student/", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "│       └── .bashrc", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "├── tmp/", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "│   └── workspace/", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "├── .ghost_dir/", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "├── secrets.txt", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "└── logs/", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "    ├── session.log", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
		{name: "    └── errors.log", offset: -30.0, velocity: 0.0, spring: harmonica.NewSpring(harmonica.FPS(60), 8.0, 0.4)},
	}

	if skipAnimations {
		for i := range files {
			files[i].revealed = true
		}
		return FileTreeModel{
			nodes:       files,
			currentFile: len(files),
			complete:    true,
			width:       width,
		}
	}

	return FileTreeModel{
		nodes:       files,
		currentFile: 0,
		complete:    false,
		width:       width,
	}
}

func (m FileTreeModel) Init() tea.Cmd {
	if m.complete {
		return nil
	}
	return tickForFileReveal()
}

func (m FileTreeModel) Update(msg tea.Msg) (FileTreeModel, tea.Cmd) {
	switch msg.(type) {
	case fileRevealMsg:
		if m.currentFile < len(m.nodes) {
			m.nodes[m.currentFile].revealed = true
			m.currentFile++
			if m.currentFile >= len(m.nodes) {
				m.complete = true
				return m, nil
			}
			return m, tickForFileReveal()
		}
	}
	return m, nil
}

func (m FileTreeModel) View() string {
	title := PanelTitleStyle(ColorBlue).Render("SANDBOX STRUCTURE")
	centeredTitle := lipgloss.Place(
		m.width-4,
		1,
		lipgloss.Center,
		lipgloss.Center,
		title,
	)

	var fileList strings.Builder
	for _, node := range m.nodes {
		if node.revealed {
			fileList.WriteString(FileTreeStyle().Render(node.name))
			fileList.WriteString("\n")
		}
	}

	return centeredTitle + "\n\n" + fileList.String()
}

func (m FileTreeModel) IsComplete() bool {
	return m.complete
}

func tickForFileReveal() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return fileRevealMsg(0)
	})
}
