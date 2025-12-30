# CLAUDE.md

This file provides guidance to Claude Code when working with this repository.

## Project Overview

RootCamp is a CLI application for learning terminal commands through interactive lessons and sandboxed labs. Built with Go, Bubble Tea (TUI), and SQLite.

## CRITICAL: /poc Directory

**DO NOT BUILD OR MODIFY CODE IN THE /poc DIRECTORY**

The `/poc` directory is for reference only. All active development and implementations must be done in the main project structure:
- Main code: `/cmd/rootcamp/` and `/internal/`
- NOT in: `/poc/` (reference only)

## Common Commands

```bash
# Build the application
go build -o rootcamp ./cmd/rootcamp

# Run the application
go run ./cmd/rootcamp

# Run with race detection (development)
go run -race ./cmd/rootcamp

# Tidy dependencies
go mod tidy

# Format code
gofmt -w .
```

## Project Structure

```
rootcamp/
├── cmd/rootcamp/          # Application entry point
├── internal/
│   ├── db/                # SQLite database layer
│   ├── lab/               # Sandbox creation/cleanup
│   ├── lessons/           # Embedded lesson definitions
│   ├── tui/               # Bubble Tea UI components
│   └── types/             # Shared data structures
├── go.mod
├── go.sum
├── PLAN.md                # Implementation plan
└── CLAUDE.md              # This file
```

## Architecture

- **TUI Framework**: Bubble Tea with Lip Gloss styling
- **Database**: SQLite stored at `~/.rootcamp/rootcamp.db`
- **Lab Sandbox**: Created at `/tmp/rootcamp-{uuid}/`, auto-cleaned on lesson exit

## Key Patterns

### Lesson Structure
Lessons are embedded Go structs in `internal/lessons/lessons.go`. Each lesson contains:
- ID (semantic: `cd`, `ls`, etc.)
- Content (What/History/CommonUses)
- Lab config (directories, files, task prompt)
- Secret code for validation

### TUI Views
- Dashboard: Lesson list with completion status
- Lesson: Content display + code input

### Progress Tracking
SQLite tracks `lesson_id`, `completed`, `completed_at`, `attempts`

## Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/charmbracelet/bubbles` - UI components
- `github.com/mattn/go-sqlite3` - SQLite driver (requires CGO)
- `github.com/google/uuid` - UUID generation for sandbox paths

## CGO Note

This project uses `go-sqlite3` which requires CGO. Ensure `CGO_ENABLED=1` and a C compiler is available:

```bash
# Check CGO status
go env CGO_ENABLED

# Build with CGO explicitly enabled
CGO_ENABLED=1 go build -o rootcamp ./cmd/rootcamp
```
