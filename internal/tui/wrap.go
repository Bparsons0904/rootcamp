package tui

import (
	"strings"
	"unicode"
)

// wrapText wraps text to the specified width, preserving existing line breaks
func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	lines := strings.Split(text, "\n")
	var wrappedLines []string

	for _, line := range lines {
		if len(line) <= width {
			wrappedLines = append(wrappedLines, line)
			continue
		}

		// Wrap this line
		wrapped := wrapLine(line, width)
		wrappedLines = append(wrappedLines, wrapped...)
	}

	return strings.Join(wrappedLines, "\n")
}

func wrapLine(line string, width int) []string {
	var result []string
	var currentLine strings.Builder
	currentLength := 0

	words := strings.Fields(line)

	for i, word := range words {
		wordLen := len(word)

		// If adding this word would exceed width
		if currentLength+wordLen+1 > width && currentLength > 0 {
			result = append(result, currentLine.String())
			currentLine.Reset()
			currentLength = 0
		}

		// Add word to current line
		if currentLength > 0 {
			currentLine.WriteString(" ")
			currentLength++
		}
		currentLine.WriteString(word)
		currentLength += wordLen

		// If this is the last word, add the line
		if i == len(words)-1 {
			result = append(result, currentLine.String())
		}
	}

	// Handle case where line was empty or only whitespace
	if len(result) == 0 && len(strings.TrimSpace(line)) == 0 {
		result = append(result, "")
	}

	return result
}

// truncate cuts text at width and adds ellipsis if needed
func truncate(text string, width int) string {
	if len(text) <= width {
		return text
	}

	if width < 3 {
		return text[:width]
	}

	// Find last space before width
	truncated := text[:width-3]
	if idx := strings.LastIndexFunc(truncated, unicode.IsSpace); idx > 0 {
		truncated = truncated[:idx]
	}

	return truncated + "..."
}
