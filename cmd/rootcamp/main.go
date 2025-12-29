package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Bparsons0904/rootcamp/internal/db"
	"github.com/Bparsons0904/rootcamp/internal/lessons"
	"github.com/Bparsons0904/rootcamp/internal/tui"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	allLessons := lessons.GetAll()

	progress, err := db.GetAllProgress(database)
	if err != nil {
		fmt.Printf("Failed to load progress: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(
		tui.NewModel(database, allLessons, progress),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
