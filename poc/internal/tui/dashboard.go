package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderDashboard() string {
	var b strings.Builder

	b.WriteString(appTitleStyle.Render(" RootCamp ") + "\n\n")
	b.WriteString(subtitleStyle.Render("Learn terminal commands through interactive lessons") + "\n\n")

	for i, lesson := range m.lessons {
		var indicator string
		var style lipgloss.Style

		progress, exists := m.progress[lesson.ID]
		if exists && progress.Completed {
			indicator = completedStyle.Render("[✓]")
		} else {
			indicator = incompleteStyle.Render("[ ]")
		}

		if i == m.selected {
			style = selectedItemStyle
			lessonText := fmt.Sprintf("%s %s - %s", indicator, lesson.ID, lesson.Title)
			b.WriteString(style.Render("▶ " + lessonText) + "\n")
		} else {
			style = unselectedItemStyle
			lessonText := fmt.Sprintf("%s %s - %s", indicator, lesson.ID, lesson.Title)
			b.WriteString(style.Render(lessonText) + "\n")
		}
	}

	b.WriteString("\n" + helpStyle.Render("↑/↓: Navigate • Enter: Start Lesson • q: Quit") + "\n")

	return b.String()
}
