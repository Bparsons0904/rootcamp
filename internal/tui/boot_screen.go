package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var systemChecks = []string{
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

type bootCheckMsg int

type BootScreenModel struct {
	currentCheck int
	complete     bool
}

func NewBootScreenModel() BootScreenModel {
	return BootScreenModel{
		currentCheck: 0,
		complete:     false,
	}
}

func (m BootScreenModel) Init() tea.Cmd {
	return tickForBootCheck()
}

func (m BootScreenModel) Update(msg tea.Msg) (BootScreenModel, tea.Cmd) {
	switch msg.(type) {
	case bootCheckMsg:
		if m.currentCheck < len(systemChecks) {
			m.currentCheck++
			if m.currentCheck >= len(systemChecks) {
				m.complete = true
				return m, nil
			}
			return m, tickForBootCheck()
		}
	}
	return m, nil
}

func (m BootScreenModel) View() string {
	var output strings.Builder
	output.WriteString("\n")

	for i := 0; i < m.currentCheck && i < len(systemChecks); i++ {
		message := systemChecks[i]

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

func (m BootScreenModel) IsComplete() bool {
	return m.complete
}

func tickForBootCheck() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return bootCheckMsg(0)
	})
}
