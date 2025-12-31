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

- **TUI Framework**: Bubble Tea with Lip Gloss styling and Huh forms
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
- Settings: Interactive form for user preferences

### Settings & Database
Settings are persisted to SQLite and loaded on startup:
- Database functions in `internal/db/db.go` handle CRUD operations
- Settings struct in `internal/types/types.go` defines available settings
- Settings UI uses Huh forms for interactive multi-select interfaces

### Bubble Tea Model Patterns

**IMPORTANT: Pointer Receivers for Forms**
When using Huh forms (or any component that binds to model fields via pointers):
- Use **pointer receivers** for the `Update()` method: `func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd)`
- Store form-bound models as **pointers** in parent models: `settingsModel *SettingsModel`
- This prevents Bubble Tea's value-based model copying from breaking pointer bindings

**Example:**
```go
// Correct - pointer receiver prevents stale pointer issues
func (m *SettingsModel) Update(msg tea.Msg) (*SettingsModel, tea.Cmd) {
    // Form binding to &m.selectedSettings works correctly
}

// Parent model stores as pointer
type Welcome3Model struct {
    settingsModel *SettingsModel  // Pointer, not value
}
```

### Progress Tracking
SQLite tracks `lesson_id`, `completed`, `completed_at`, `attempts`

## Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/charmbracelet/bubbles` - UI components
- `github.com/charmbracelet/huh` - Interactive forms (PREFERRED for forms/inputs)
- `github.com/charmbracelet/glamour` - Markdown rendering
- `github.com/charmbracelet/harmonica` - Spring animations
- `github.com/mattn/go-sqlite3` - SQLite driver (requires CGO)
- `github.com/google/uuid` - UUID generation for sandbox paths

**Preference: Use Huh for Forms**
When implementing any user input forms (settings, configuration, multi-select, etc.), prefer using the `huh` library. It provides:
- Built-in validation
- Consistent styling and theming
- Multi-select, input, confirm, select, and note components
- Better user experience than manual input handling

## Development Practices

### Code Organization
- Keep comments minimal and purposeful - code should be self-documenting
- Use template comments for future feature additions (commented-out examples)
- Remove experimental/example code before finalizing features

### Settings Implementation Pattern
When adding new settings:
1. Add field to `types.Settings` struct
2. Add default value in `db.InitDefaultSettings()`
3. Add option in `SettingsModel.Open()` function
4. Add loading logic in `db.GetAllSettings()`
5. Apply setting where needed (e.g., in model initialization)

### Animation & UX Patterns
- Use `skippedAnimations` flag to conditionally render progress indicators
- Check settings during model initialization, not every render
- Use phase-based state management for multi-step animations

### Embedded Content: Multi-File Loading Pattern

**Architecture:** Fun facts are loaded from multiple JSON files using Go's `embed.FS`

**Pattern:** `internal/lessons/funfacts.go`
```go
//go:embed data/funfacts/*.json
var embeddedFunFactsFS embed.FS

func LoadFunFacts() (*types.FunFactsData, error) {
    entries, err := embeddedFunFactsFS.ReadDir("data/funfacts")
    // Iterate through all .json files
    // Merge facts from each file into single array
}
```

**Benefits:**
- Add/remove JSON files without modifying Go code
- Organize content by category (commands.json, terminal.json, etc.)
- All files merged into single cached data structure
- Standard library approach (no external dependencies)

**Directory Structure:**
```
internal/lessons/data/funfacts/
├── all.json           # Original facts
├── modern-tools.json  # Modern CLI tools
└── <category>.json    # Additional themed files
```

### TUI Performance: Pre-Rendering Pattern

**CRITICAL for TUI responsiveness:** Never perform expensive operations (markdown rendering, heavy parsing) on user interaction.

**Problem:** On-demand Glamour markdown rendering causes multi-second delays when selecting facts.

**Solution:** Pre-render during model initialization, cache in map for O(1) lookup.

**Pattern:** `internal/tui/fun_facts.go`
```go
type FunFactsModel struct {
    allFacts      []types.FunFact
    renderedFacts map[string]string  // Pre-rendered markdown cache
}

func NewFunFactsModel(database *sql.DB) FunFactsModel {
    // Load facts
    data, _ := lessons.LoadFunFacts()

    // Pre-render ALL markdown at startup
    renderedFacts := make(map[string]string)
    renderer, _ := glamour.NewTermRenderer(...)

    for _, fact := range data.Facts {
        rendered, _ := renderer.Render(fact.Full)
        renderedFacts[fact.ID] = rendered  // Cache it
    }

    return FunFactsModel{
        allFacts:      data.Facts,
        renderedFacts: renderedFacts,
    }
}

func (m *FunFactsModel) setupDetailView() {
    // Instant lookup - no rendering!
    rendered := m.renderedFacts[m.selectedFactID]
    m.viewport.SetContent(rendered)
}
```

**Result:**
- Startup: All rendering happens once (acceptable delay)
- User interaction: Instant (simple map lookup)
- TUI stays snappy and responsive

**General Rule:** For TUI apps, prefer:
- Pre-computation over on-demand computation
- Map lookups over repeated operations
- Cached data over live generation

## CGO Note

This project uses `go-sqlite3` which requires CGO. Ensure `CGO_ENABLED=1` and a C compiler is available:

```bash
# Check CGO status
go env CGO_ENABLED

# Build with CGO explicitly enabled
CGO_ENABLED=1 go build -o rootcamp ./cmd/rootcamp
```
