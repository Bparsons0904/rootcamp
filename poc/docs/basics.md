🏟️ Project Manifest: RootCamp (ttry)RootCamp is an open-source, terminal-native training platform designed to move engineers from GUI-dependency to CLI-proficiency. It uses a "sandbox and secret" loop to provide hands-on, high-fidelity practice.🎯 Core PhilosophyCLI-First: The tool should feel like a native Unix utility (short name, fast, pipe-friendly).High-Fidelity Labs: Users don't just "read" about commands; they execute them in a real filesystem sandbox.Mentor Vibe: Lessons provide context (etymology/history) to make the knowledge "stick."Low Friction: No complex setup. A single binary + a local SQLite DB.🛠️ Technical StackComponentTechnologyRationaleLanguageGoPerformance, static typing, and easy cross-platform binary distribution.TUI FrameworkBubble Tea (Charm.sh)The gold standard for modern, beautiful, and reactive terminal interfaces.DatabaseSQLiteZero-config, single-file persistence for progress tracking.SandboxingTemporary DirectoriesUtilizes /tmp/rootcamp-{uuid} to ensure safety and isolation.🏗️ Data ArchitectureWe will use semantic IDs (e.g., cd-01) to allow for easy grouping and expansion.Go Struct DefinitionsGotype Lesson struct {
ID string // e.g., "nav-cd-01"
Title string // e.g., "cd"
Category string // "navigation", "manipulation", "permissions"
What string // Brief definition
Example string // Syntax example
History string // Why it exists (The "Mentor" context)
CommonUses []string // Practical shortcuts/flags
LabSetup LabConfig // Instructions for generating the sandbox
SecretCode string // The validation key
Difficulty int // 1-5
}

type LabConfig struct {
Files map[string]string // Path -> Content (e.g., "README.md": "Code: 123")
Dirs []string // Directories to create
Task string // The specific instruction for the user
}

type Progress struct {
LessonID string
Completed bool
Attempts int
LastSeen time.Time
}
🔄 The "Lab Loop" WorkflowSelection: User chooses a lesson via the TUI dashboard.Provisioning: RootCamp generates a random UUID and populates /tmp/rootcamp-{uuid} based on the LabConfig.The Prompt: The TUI displays the "How/Why" and the "Task."The Action: The user opens a second terminal tab (or uses a multiplexer like tmux/zellij) to perform the task.Validation: User finds the "Secret Code" in the sandbox and enters it into the RootCamp TUI.Teardown: On success (or exit), the sandbox is scrubbed.🗺️ v0.1 Roadmap: The "Core Five"We will launch with five foundational commands to prove the concept:cd: Navigating nested structures.ls: Discovering hidden files and permissions (e.g., ls -la).mv: Renaming vs. moving files.cp: Recursive copying of directories.cat: Reading file contents to find the secret.🏁 Design Decisions (Resolved)Naming: We will use RootCamp as the project name, but the binary will be ttry for quick invocation.IDs: Semantic naming (e.g., filesystem-ls-01) to allow for modular lesson packs.Cleanup: Auto-delete /tmp sandboxes upon TUI exit to keep the user's system clean.Progression: A "Recommended Path" (Ordered) will be the default view, though all lessons are unlocked by default for "Power Users."
