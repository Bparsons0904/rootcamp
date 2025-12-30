This updated capability matrix reflects the state of the **Charm** ecosystem as of **late 2025**, ensuring we leverage the most current features for `ttry`.

---

## 🛠️ RootCamp Tooling Capability Matrix (2025 Edition)

This document serves as our technical reference for what is possible within the terminal interface.

| Tool           | Core Functionality     | 2025 High-Fidelity Features                                                                                                                                                                                                                       |
| -------------- | ---------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Bubble Tea** | **TUI Engine**         | High-performance framerate renderer, built-in **Focus/Blur** reporting (vital for detecting when a user switches back from their lab terminal), and `BatchMsg` for handling concurrent I/O (e.g., watching the filesystem while updating the UI). |
| **Lip Gloss**  | **Layout & Style**     | New **Layer/Canvas** system for overlapping elements (popups/modals), adaptive color profiles (4-bit to True Color), and CSS-like margin/padding shorthand.                                                                                       |
| **Huh?**       | **Interactive Forms**  | Group-based "multi-page" forms for sequential tasks, real-time input validation, and an accessible mode for screen readers.                                                                                                                       |
| **Bubbles**    | **UI Components**      | A massive library including **Fuzzy-filtering Lists**, responsive **Viewports**, high-performance **Tables** (with wrapping support), and specialized widgets like **Filepickers** and **Stopwatches**.                                           |
| **Glamour**    | **Markdown Rendering** | Stylesheet-driven rendering with **Chroma syntax highlighting**, automatic word-wrapping, and the new "Tokyo Night" theme.                                                                                                                        |
| **Harmonica**  | **Animation**          | Physics-based spring animations to make progress bars and transitions feel fluid rather than "jumpy".                                                                                                                                             |

---

## 🎭 The RootCamp Persona: "The Architect"

To make the knowledge "stick" as per our core philosophy, we need a persona that commands respect but remains an accessible peer.

**Name:** The Lead Architect

**Vibe:** Experienced, precise, and slightly "old-school" cool. Think of a senior engineer who has been using the command line since the 90s but loves modern, clean design.

**Tone:**

- **Encouraging but Direct:** "You're leaning on the mouse again. Let's fix that."
- **Historical Context:** Connects a simple `cd` command to the physical architecture of early disk drives (The Mentor Vibe).
- **Empathetic:** Acknowledges that the CLI is intimidating but reminds the user that it’s where "real power lives."
- **Intellectually Honest:** If a user makes a mistake, the persona won't just say "Error"; it will explain _why_ the filesystem rejected the move (e.g., "Permissions aren't just suggestions; they're the law of the kernel").

---

## 🏗️ Design Concept: The "Nerve Center" Dashboard

By combining **Lip Gloss**'s new layer system with **Bubbles**' list component, we can create a dashboard that feels like a mission control center.

1. **Header:** A bold, centered title using `lipgloss.NewStyle().Border(lipgloss.DoubleBorder())`.
2. **Navigation (Left):** A `list.Model` that lets users filter through the "Core Five" labs.
3. **Details (Right):** A `viewport.Model` rendering the "Mentor History" via **Glamour**.
4. **Footer:** A `progress.Model` and a `help.Model` that dynamically updates based on whether a lab is "Provisioning," "In-Progress," or "Validated."

**Would you like me to write the Go code for the `MainDashboard` model using the Bubble Tea architecture we just defined?**
