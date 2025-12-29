package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderLesson() string {
	if m.currentLesson == nil {
		return "No lesson selected"
	}

	lesson := *m.currentLesson
	var b strings.Builder

	contentWidth := m.width - 4
	if contentWidth < 40 {
		contentWidth = 40
	}

	b.WriteString(appTitleStyle.Render(" RootCamp ") + "\n")
	b.WriteString(titleStyle.Render(fmt.Sprintf("%s - %s", lesson.ID, lesson.Title)) + "\n\n")

	b.WriteString(sectionTitleStyle.Render("What:") + " ")
	b.WriteString(textStyle.Render(wrapText(lesson.What, contentWidth-6)) + "\n\n")

	b.WriteString(sectionTitleStyle.Render("Task:") + " ")
	b.WriteString(taskStyle.Render(wrapText(lesson.Lab.Task, contentWidth-6)) + "\n\n")

	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444")).
		Render(strings.Repeat("─", contentWidth))
	b.WriteString(divider + "\n\n")

	if m.codeInputActive {
		b.WriteString(inputStyle.Render("Secret Code: ") + m.input.View() + "\n")
		if m.feedback != "" {
			if m.feedbackSuccess {
				b.WriteString(successStyle.Render(m.feedback) + "\n")
			} else {
				b.WriteString(errorStyle.Render(m.feedback) + "\n")
			}
		}
		b.WriteString("\n" + helpStyle.Render("Enter: Submit • Esc: Cancel") + "\n")
	} else if m.labStarted {
		statusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)
		b.WriteString(statusStyle.Render("Lab completed!") + "\n\n")
		b.WriteString(textStyle.Render("You've practiced in the terminal. Now enter the secret code you found.") + "\n\n")
		b.WriteString(helpStyle.Render("Enter: Start lab again • Ctrl+D: Enter secret code • Esc: Back to dashboard") + "\n")
	} else if m.activeLab != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true)
		b.WriteString(statusStyle.Render("Lab in progress...") + "\n\n")
		b.WriteString(helpStyle.Render("Esc: Back to dashboard") + "\n")
	} else {
		actionStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)
		b.WriteString(actionStyle.Render("Press Enter to start the lab") + "\n\n")
		b.WriteString(textStyle.Render("You'll be dropped into a terminal in the sandbox environment.") + "\n")
		b.WriteString(textStyle.Render("Type 'exit' when done to return here.") + "\n\n")
		b.WriteString(helpStyle.Render("Enter: Start lab • Ctrl+D: Enter secret code • Esc: Back to dashboard") + "\n")
	}

	return b.String()
}
