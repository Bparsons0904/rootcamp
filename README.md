# RootCamp

A CLI application for learning terminal commands through interactive lessons and sandboxed labs.

## Features

- **Interactive TUI** - Beautiful terminal user interface built with Bubble Tea
- **Sandboxed Labs** - Practice commands in isolated environments
- **Progress Tracking** - SQLite database tracks your completion status
- **5 Starter Lessons** - Learn cd, ls, cat, cp, and mv commands

## Installation

### Prerequisites

- Go 1.21 or later
- CGO-enabled (required for SQLite)
- C compiler (gcc, clang, etc.)

### Building

```bash
# Clone the repository
git clone https://github.com/Bparsons0904/rootcamp.git
cd rootcamp

# Build the application
CGO_ENABLED=1 go build -o rootcamp ./cmd/rootcamp

# Run the application
./rootcamp
```

## Usage

### Dashboard Navigation

- **↑/↓ or k/j** - Navigate between lessons
- **Enter** - Start selected lesson
- **Ctrl+C** - Quit application

### In-Lesson Controls

RootCamp features an **embedded terminal** that runs directly in the lesson view!

#### Terminal Mode (default)
- **Type normally** - All keystrokes go to the embedded shell
- **Ctrl+D** - Switch to code input mode
- **Esc** - Return to dashboard (closes terminal and sandbox)

#### Code Input Mode
- **Type your code** - Enter the secret code you discovered
- **Enter** - Submit code for validation
- **Esc** - Return to terminal mode

### Completing Lessons

1. Select a lesson from the dashboard
2. Read the lesson content and task at the top
3. Use the **embedded terminal** in the bottom half to complete the task
4. Find the secret code in the sandbox
5. Press **Ctrl+D** to enter code input mode
6. Type the secret code and press **Enter**
7. Move on to the next lesson!

## Lessons

1. **cd** - Navigate directories
2. **ls** - List directory contents
3. **cat** - Display file contents
4. **cp** - Copy files and directories
5. **mv** - Move and rename files

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
└── README.md
```

## Technical Details

- **TUI Framework**: Bubble Tea with Lip Gloss styling
- **Embedded Terminal**: PTY-based shell integration for in-app command practice
- **Database**: SQLite (`~/.rootcamp/rootcamp.db`)
- **Lab Sandbox**: Temporary directories at `/tmp/rootcamp-{uuid}/`
- **Auto-cleanup**: Labs and terminals are automatically cleaned up when exiting lessons
- **Split-Screen Layout**: Lesson content on top, interactive terminal on bottom

## License

MIT License
