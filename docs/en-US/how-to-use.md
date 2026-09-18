# How to use tuitodo

**Language:** English · [Portuguese (Brazil)](../pt-BR/how-to-use.md)

## Run

Requires Go 1.25 or newer.

```bash
go run ./cmd/tuitodo
```

## Layout

Four panes:

- **Tasks** — modes **To do** and **Completed**
- **Categories** — includes the virtual category **All**
- **Subtasks** — children of the selected parent (one level)
- **Actions** — hints for the current pane

## Navigation

- Arrow keys choose the pane
- **Enter** enters the pane
- **Esc** goes back

Keyboard shortcuts stay the same in every language. Pane titles and labels follow the selected language. Task titles are your data.

## Shortcuts

- `1` — To do
- `2` — Completed
- `a` — new task (in Subtasks: new subtask)
- `c` — new category
- `s` — Settings
- `e` — edit
- `d` — delete
- `Space` — complete (in Completed: reopen)
- `?` — shortcuts help
- `q` — quit

## Completed

The left list shows **dates only** (the local day of the parent's completion). The right pane lists that day's parent tasks and their subtasks. **Space** reopens the focused item. Reopening a subtask returns the family to To do.

## Language

On first launch the app follows the OS (`en-US` or `pt-BR`, falling back to `en-US`). Press `s` to open Settings. Changing the language in the list applies immediately. **Save** stores it; **Cancel** or **Esc** reverts.

## Data

Tasks live in a local SQLite file:

- macOS: `~/Library/Application Support/tuitodo/tuitodo.db`
- Linux: `~/.local/share/tuitodo/tuitodo.db` (or `$XDG_DATA_HOME/tuitodo/tuitodo.db`)
