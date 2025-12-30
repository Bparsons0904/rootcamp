# RootCamp v0.1 POC Implementation Plan

## Overview
CLI app for learning terminal commands using Bubble Tea TUI, SQLite persistence, and sandboxed labs.

## Decisions Made
- **Lesson IDs**: Semantic (`cd`, `ls`, `mv`, `cp`, `cat`)
- **Lab cleanup**: Auto-delete on lesson completion
- **Attempts**: Track in DB, don't display in v0.1
- **Lesson order**: Curated (manual sequence)
- **Module**: `github.com/Bparsons0904/rootcamp`
- **Lesson storage**: Embedded Go structs
- **Structure**: Standard Go layout (`cmd/`, `internal/`)

---

## Project Structure

```
rootcamp/
├── cmd/
│   └── rootcamp/
│       └── main.go           # Entry point, TUI initialization
├── internal/
│   ├── db/
│   │   └── db.go             # SQLite setup, queries
│   ├── lab/
│   │   └── lab.go            # Sandbox creation/cleanup
│   ├── lessons/
│   │   └── lessons.go        # Embedded lesson definitions
│   ├── tui/
│   │   ├── app.go            # Main Bubble Tea model
│   │   ├── dashboard.go      # Lesson list view
│   │   ├── lesson.go         # Lesson detail view
│   │   └── styles.go         # Lip Gloss styling
│   └── types/
│       └── types.go          # Lesson, LabConfig, UserProgress structs
├── go.mod
├── go.sum
└── README.md
```

---

## Implementation Steps

### Step 1: Project Initialization
1. `go mod init github.com/Bparsons0904/rootcamp`
2. Create directory structure
3. Add dependencies:
   - `github.com/charmbracelet/bubbletea`
   - `github.com/charmbracelet/lipgloss`
   - `github.com/charmbracelet/bubbles`
   - `github.com/mattn/go-sqlite3`
   - `github.com/google/uuid`

### Step 2: Types (`internal/types/types.go`)
```go
type Lesson struct {
    ID          string
    Title       string
    Order       int           // Curated sequence position
    What        string
    Example     string
    History     string
    CommonUses  []string
    Lab         LabConfig
    SecretCode  string
    Difficulty  string
    Tags        []string
}

type LabConfig struct {
    Dirs   []string            // Directories to create
    Files  map[string]string   // filepath -> content
    Task   string              // User prompt
}

type UserProgress struct {
    LessonID    string
    Completed   bool
    CompletedAt *time.Time
    Attempts    int
}
```

### Step 3: Database (`internal/db/db.go`)
- Create `~/.rootcamp/rootcamp.db`
- Schema:
  ```sql
  CREATE TABLE IF NOT EXISTS progress (
      lesson_id TEXT PRIMARY KEY,
      completed BOOLEAN DEFAULT FALSE,
      completed_at DATETIME,
      attempts INTEGER DEFAULT 0
  );
  ```
- Functions:
  - `InitDB() (*sql.DB, error)`
  - `GetProgress(db, lessonID) (*UserProgress, error)`
  - `MarkComplete(db, lessonID) error`
  - `IncrementAttempts(db, lessonID) error`

### Step 4: Lessons (`internal/lessons/lessons.go`)
Define 5 starter lessons as embedded Go structs:

1. **cd** (Order: 1) - Navigate directories
2. **ls** (Order: 2) - List directory contents
3. **cat** (Order: 3) - Display file contents
4. **cp** (Order: 4) - Copy files
5. **mv** (Order: 5) - Move/rename files

Each lesson includes:
- What it does + example
- History/etymology
- 5 common use cases
- Lab config (dirs, files, task)
- Secret code

### Step 5: Lab Manager (`internal/lab/lab.go`)
- `Create(lesson Lesson) (sandboxPath string, error)`
  - Creates `/tmp/rootcamp-{uuid}/`
  - Builds directory structure from `lesson.Lab.Dirs`
  - Writes files from `lesson.Lab.Files`
  - Returns sandbox path
- `Cleanup(sandboxPath string) error`
  - Removes sandbox directory
- `Validate(userInput, secretCode string) bool`
  - Compares user input to expected code

### Step 6: TUI (`internal/tui/`)

**app.go** - Main model:
```go
type Model struct {
    db          *sql.DB
    lessons     []Lesson
    progress    map[string]*UserProgress
    view        View  // Dashboard or Lesson
    selected    int
    activeLab   string
    input       textinput.Model
    // ...
}
```

**dashboard.go** - Lesson list:
- Shows all 5 lessons with completion status
- `[✓] cd - Navigate directories`
- `[ ] ls - List contents`
- Up/Down to navigate, Enter to select

**lesson.go** - Lesson view:
- Displays What/History/Common Uses
- Shows task prompt
- Text input for secret code
- Feedback on submit (correct/incorrect)
- Auto-cleanup lab on exit

**styles.go** - Lip Gloss styles:
- Title styling
- Completed vs incomplete indicators
- Input field styling
- Feedback colors (green success, red error)

### Step 7: Entry Point (`cmd/rootcamp/main.go`)
```go
func main() {
    db := db.InitDB()
    defer db.Close()

    lessons := lessons.GetAll()
    progress := db.GetAllProgress()

    p := tea.NewProgram(tui.NewModel(db, lessons, progress))
    p.Run()
}
```

---

## File Creation Order

1. `go.mod` - Module init
2. `internal/types/types.go` - Data structures
3. `internal/lessons/lessons.go` - 5 hardcoded lessons
4. `internal/db/db.go` - SQLite setup + queries
5. `internal/lab/lab.go` - Sandbox manager
6. `internal/tui/styles.go` - Styling
7. `internal/tui/dashboard.go` - Lesson list view
8. `internal/tui/lesson.go` - Lesson detail view
9. `internal/tui/app.go` - Main TUI model
10. `cmd/rootcamp/main.go` - Entry point

---

## POC Scope (Explicitly NOT Included)

- Streaks/gamification
- XP/leveling
- Multiple difficulty levels
- Hints system
- Lesson search/filtering
- External lesson files
- Tests (minimal implementation per user preference)
