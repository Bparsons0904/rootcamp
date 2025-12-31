package tui

import (
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var architectFacts = []string{
	`In **1965**, Multics introduced hierarchical directories.

Before this, data was a _flat pile of magnetic tape_—no hierarchy.

The directory changed everything:
- **Namespaces** for organization
- **Paths** to navigate data
- **Permissions** to control access

Every command you run traverses this tree.`,

	`The **pipe** (|) was invented by Doug McIlroy in **1973**.

His vision: _"Write programs that do one thing well and work together."_

This single character changed software:
- **Composability** of tools
- The Unix **philosophy**
- Simple power: **cat log | grep ERROR**`,

	`In **1991**, Linus Torvalds posted to comp.os.minix:

_"I'm doing a (free) OS (just a hobby)..."_

That hobby became **Linux**:
- Powers **96.3%** of top servers
- Every **Android** device
- **100%** of top 500 supercomputers

A student project, now foundation of the internet.`,

	`**/bin** and **/usr/bin** split: **1971**.

Unix ran out of disk space. The PDP-11 had a **1.5MB** drive. Dennis Ritchie added a second disk at /usr.

Today:
- **/bin** - System binaries
- **/usr/bin** - User programs

Modern Linux carries the ghost of a 1970s disk shortage.`,

	`**chmod** uses octal due to **1974** hardware limits.

Permissions needed **9 bits** (rwxrwxrwx). Octal aligned perfectly:
- **755** = rwxr-xr-x
- **644** = rw-r--r--

Memory was expensive, every bit counted.

We still use octal because _that's how it's always been_.`,

	`The **root** user was never meant to be permanent.

Ken Thompson created it for testing. It was **temporary**.

Instead, it became **immortal**:
- Every Unix system has root
- 50+ years later, still here

The ultimate _"temporary solution"_.`,

	`**Hidden files** (.bashrc) were an accident.

Early **ls** sorted alphabetically. Files starting with **.** sorted first.

Later, someone hid dotfiles to reduce clutter.

Result:
- Configs became "special"
- Pattern became **convention**

Your home is littered with dotfiles from a sorting hack.`,

	`**/etc** means **"et cetera"**—_"and other things."_

Early Unix had:
- **/bin** for binaries
- **/dev** for devices
- **/lib** for libraries

Everything else? **Et cetera.**

The "misc folder" became the backbone of system admin.`,

	`**tty** = **teletypewriter** (1920s hardware).

Early terminals were literal typewriters:
- No screen, just paper
- Type a command, it prints

Modern terminals are **emulators** of 100-year-old machines.

A simulation of a simulation.`,

	`The **$** prompt has military origins.

In the **1960s**, computing cost money per CPU second. **$** reminded users:
_"This costs money"_

Root used **#** (override costs).

Today, free computing, but **$** and **#** remain.`,
}

type ArchitectLogModel struct {
	selectedFact    string
	glamourRenderer *glamour.TermRenderer
	width           int
}

func NewArchitectLogModel(width int) ArchitectLogModel {
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(40),
	)

	randomIndex := rand.Intn(len(architectFacts))

	return ArchitectLogModel{
		selectedFact:    architectFacts[randomIndex],
		glamourRenderer: renderer,
		width:           width,
	}
}

func (m ArchitectLogModel) Init() tea.Cmd {
	return nil
}

func (m ArchitectLogModel) Update(msg tea.Msg) (ArchitectLogModel, tea.Cmd) {
	return m, nil
}

func (m ArchitectLogModel) View() string {
	title := PanelTitleStyle(ColorPurple).Render("THE ARCHITECT'S LOG")
	centeredTitle := lipgloss.Place(
		m.width-4,
		1,
		lipgloss.Center,
		lipgloss.Center,
		title,
	)

	rendered, err := m.glamourRenderer.Render(m.selectedFact)
	if err != nil {
		return centeredTitle + "\n\n" + m.selectedFact
	}

	return centeredTitle + "\n\n" + strings.TrimSpace(rendered)
}
