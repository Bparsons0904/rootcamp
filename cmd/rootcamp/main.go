package main

import (
	"fmt"
	"os"

	"github.com/bobparsons/rootcamp/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := tui.NewWelcome3Model()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
