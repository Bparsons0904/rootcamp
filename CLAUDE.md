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
- **Lab Sandbox**: Created at `/tmp/rootcamp-{5char}/`, auto-cleaned on lesson exit (e.g., `/tmp/rootcamp-a3x9z`)

### Learn Command vs Guided Learning

**IMPORTANT DISTINCTION:**
- **Learn Command** (IMPLEMENTED): Individual lesson selection - user picks any lesson and completes it independently
- **Guided Learning** (FUTURE): Predetermined learning path - system decides lesson order based on progression

These are separate features. "Learn Command" is the menu item for the implemented config-driven lesson system.

## Key Patterns

### Lesson Structure

**Config-Driven Lessons:** Lessons are defined in JSON files following the multi-file loading pattern.

**Location:** `internal/lessons/data/lessons/*.json`

**Structure:** Each lesson contains:
- **ID/Code**: Semantic identifier (`pwd`, `ls`, `cd`, etc.)
- **About**: Educational content (What, History, Example, CommonUses)
- **Sandbox**: Lab environment config (directories, files, startDir)
- **Instructions**: Markdown task description (includes "## Your Task" heading)
- **Requirements**: Validation criteria (type, validator, expected value)
- **Hints**: Tips for completion (stored but NOT displayed - reserved for future feature)

**Pattern:** Follows `internal/lessons/funfacts.go` multi-file loading:
```go
//go:embed data/lessons/*.json
var embeddedLessonsFS embed.FS

func LoadLessons() (*types.LessonsData, error) {
    // Merges all JSON files into single cached structure
}
```

**Key Detail:** Instructions field contains full markdown including headings. Hints are excluded from display.

### TUI Views
- **Learn Command List**: Centered lesson selection form (90 width, centered using lipgloss.Place)
- **Lesson Detail**: Horizontally centered content (90 width, full height viewport, top-aligned)
- **Code Input**: User answer entry after exiting sandbox
- **Success**: Completion confirmation screen
- **Settings**: Interactive form for user preferences

### TUI Layout Patterns

**Centering Content:**
```go
// Center both horizontally and vertically (lesson list)
return lipgloss.Place(
    m.width,
    m.height,
    lipgloss.Center,
    lipgloss.Center,
    content,
)

// Center horizontally, top-align vertically (lesson detail)
return lipgloss.Place(
    m.width,
    m.height,
    lipgloss.Center,
    lipgloss.Top,
    content,
)
```

**Content Width Consistency:**
- Form width: 90 characters (set in createForm)
- Detail view: 90 characters (matches form for visual consistency)
- Viewport: 90 width, dynamic height (m.height - 8 for UI chrome)

**Key Principle:** Constrain content width for readability, center horizontally, use full screen height

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

### Lab/Sandbox Startup Flow

**Auto-Spawn Pattern:** Terminal automatically opens and displays instructions when user starts lab.

**Implementation:** `internal/tui/learn_command.go` - `startLab()`
```go
instructions := fmt.Sprintf(`clear
cat << 'EOF'
╔═══════════════════════════════════════════════╗
║           ROOT CAMP - LAB SESSION            ║
╚═══════════════════════════════════════════════╝

Lesson: %s
%s

Your sandbox is located at: %s
When you're done, type 'exit' to return to Root Camp.
EOF
exec bash`, lesson.Title, lesson.Instructions, startPath)

c := exec.Command("bash", "-c", instructions)
c.Env = append(os.Environ(), fmt.Sprintf("PS1=rootcamp:%s$ ", lesson.Code))
```

**Flow:**
1. User presses 'S' to start lab
2. `bash -c` executes script that:
   - Clears terminal (`clear`)
   - Prints instructions via heredoc
   - Replaces shell with interactive bash (`exec bash`)
3. User completes task in sandbox
4. User types `exit` → returns to Root Camp code input screen

**Sandbox ID Generation:**
```go
// internal/lab/sandbox.go
func generateShortID() string {
    const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
    const length = 5
    // Returns 5-char alphanumeric (e.g., "a3x9z")
}
```

**Key Benefits:**
- Clean terminal on entry
- Instructions visible while working
- Simple exit mechanism
- Readable sandbox paths (`/tmp/rootcamp-a3x9z`)

### Progress Tracking
SQLite tracks `lesson_id`, `completed`, `completed_at`, `attempts`

### Requirement Validation
Flexible validation system in `internal/lab/validate.go`:
- **exact**: Simple string match
- **path_match**: Replaces `{uuid}` placeholder with actual sandbox ID before comparison
- **file_check**: Verifies file existence in sandbox
- **regex**: Pattern matching

## Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/charmbracelet/bubbles` - UI components
- `github.com/charmbracelet/huh` - Interactive forms (PREFERRED for forms/inputs)
- `github.com/charmbracelet/glamour` - Markdown rendering
- `github.com/charmbracelet/harmonica` - Spring animations
- `github.com/mattn/go-sqlite3` - SQLite driver (requires CGO)

**Note:** Sandbox IDs use custom 5-character alphanumeric generation (stdlib `math/rand`), not UUIDs.

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

### Lesson Implementation Pattern

**When adding new lessons:**

1. **Create JSON file** in `internal/lessons/data/lessons/*.json`
2. **Structure** each lesson with:
   - Unique ID/code (command name)
   - Complete About section (what, history, example, commonUses)
   - Instructions with markdown heading: `"## Your Task\n\n..."`
   - Sandbox config (dirs, files, startDir)
   - Requirements array for validation
   - Hints array (stored but not displayed)

3. **Instructions format:**
   - Include full markdown with heading
   - Don't duplicate headings (formatLessonAbout displays instructions directly)
   - Be clear about task and exit mechanism

4. **Validation:**
   - Use `{uuid}` placeholder in expected paths (replaced with actual sandbox ID)
   - Choose appropriate validator: exact, path_match, file_check, regex

**Example structure:**
```json
{
  "instructions": "## Your Task\n\nRun the `pwd` command...",
  "requirements": [{
    "validator": "path_match",
    "expected": "/tmp/rootcamp-{uuid}/projects/rootcamp"
  }]
}
```

## CGO Note

This project uses `go-sqlite3` which requires CGO. Ensure `CGO_ENABLED=1` and a C compiler is available:

```bash
# Check CGO status
go env CGO_ENABLED

# Build with CGO explicitly enabled
CGO_ENABLED=1 go build -o rootcamp ./cmd/rootcamp
```
