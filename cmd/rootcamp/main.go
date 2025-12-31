package main

import (
	"fmt"
	"os"

	"github.com/bobparsons/rootcamp/internal/db"
	"github.com/bobparsons/rootcamp/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	model := tui.NewWelcomeModel(database)
	p := tea.NewProgram(&model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
