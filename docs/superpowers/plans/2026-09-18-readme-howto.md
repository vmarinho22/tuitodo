# README e how-to-use Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add an English README that points to complete how-to-use guides in English and Brazilian Portuguese.

**Architecture:** Three markdown files only. README is a short English entry point (title, one context paragraph, two links). Each locale has one `docs/<locale>/how-to-use.md` with the same section headings, a language switcher at the top, and copy that matches the current TUI labels and shortcuts. No code changes.

**Tech Stack:** Markdown. App facts come from `internal/i18n/locales/*/messages.json`, `internal/ui/app.go` (keys `1` `2` `a` `c` `s` `e` `d` ` ` `?` `q`), and `internal/store/sqlite.go` `DefaultDatabasePath`.

## Global Constraints

- README is English only; no `README.pt-BR.md`
- README has no install, shortcuts, badges, or development section
- How-to-use link labels: `English` and `Portuguese (Brazil)`
- Both guides share the same structure: language switcher, then Run, Layout, Navigation, Shortcuts, Completed, Language, Data
- Keyboard shortcuts do not change with language; UI labels do; task titles are user data
- `docs/superpowers/` stays internal; README and how-to-use do not link there
- No `docs/*/how-to-use/` directory of multiple files
- No screenshots

## File map

| File | Responsibility |
|------|----------------|
| `README.md` | English title, context paragraph, How to use links |
| `docs/en-US/how-to-use.md` | Full English guide |
| `docs/pt-BR/how-to-use.md` | Full Portuguese (Brazil) guide |

---

### Task 1: README.md

**Files:**
- Create: `README.md`

**Interfaces:**
- Consumes: nothing
- Produces: relative links `docs/en-US/how-to-use.md` and `docs/pt-BR/how-to-use.md` (files may not exist until later tasks)

- [x] **Step 1: Write README.md**

Create `README.md` with exactly:

```markdown
# tuitodo

Terminal todo app: parent tasks with one-level subtasks, categories,
and a local SQLite database. The UI is a TUI (tview). Language follows
the OS (en-US / pt-BR) and can be changed in the app.

## How to use

- [English](docs/en-US/how-to-use.md)
- [Portuguese (Brazil)](docs/pt-BR/how-to-use.md)
```

- [x] **Step 2: Verify README**

Run:

```bash
test -f README.md
grep -F '# tuitodo' README.md
grep -F '[English](docs/en-US/how-to-use.md)' README.md
grep -F '[Portuguese (Brazil)](docs/pt-BR/how-to-use.md)' README.md
! grep -Ei 'install|badge|go install|docs/superpowers' README.md
```

Expected: `test` and the three `grep -F` succeed; the last command exits 0 (no matches).

- [x] **Step 3: Commit**

```bash
git add README.md
git commit -m "$(cat <<'EOF'
docs: add English README with how-to-use links

EOF
)"
```

---

### Task 2: English how-to-use

**Files:**
- Create: `docs/en-US/how-to-use.md`

**Interfaces:**
- Consumes: README link `docs/en-US/how-to-use.md`
- Produces: sibling link `../pt-BR/how-to-use.md` (created in Task 3)

- [x] **Step 1: Write docs/en-US/how-to-use.md**

Create `docs/en-US/how-to-use.md` with exactly:

```markdown
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
```

- [x] **Step 2: Verify English guide**

Run:

```bash
test -f docs/en-US/how-to-use.md
grep -F '## Run' docs/en-US/how-to-use.md
grep -F '## Layout' docs/en-US/how-to-use.md
grep -F '## Navigation' docs/en-US/how-to-use.md
grep -F '## Shortcuts' docs/en-US/how-to-use.md
grep -F '## Completed' docs/en-US/how-to-use.md
grep -F '## Language' docs/en-US/how-to-use.md
grep -F '## Data' docs/en-US/how-to-use.md
grep -F 'go run ./cmd/tuitodo' docs/en-US/how-to-use.md
grep -F '[Portuguese (Brazil)](../pt-BR/how-to-use.md)' docs/en-US/how-to-use.md
! grep -F 'docs/superpowers' docs/en-US/how-to-use.md
```

Expected: all `grep -F` succeed; the last command exits 0 (no matches).

- [x] **Step 3: Commit**

```bash
git add docs/en-US/how-to-use.md
git commit -m "$(cat <<'EOF'
docs: add English how-to-use guide

EOF
)"
```

---

### Task 3: Portuguese (Brazil) how-to-use

**Files:**
- Create: `docs/pt-BR/how-to-use.md`

**Interfaces:**
- Consumes: English switcher link `../pt-BR/how-to-use.md`
- Produces: switcher link `../en-US/how-to-use.md`

- [x] **Step 1: Write docs/pt-BR/how-to-use.md**

Create `docs/pt-BR/how-to-use.md` with exactly:

```markdown
# Como usar o tuitodo

**Idioma:** [English](../en-US/how-to-use.md) · Portuguese (Brazil)

## Executar

Requer Go 1.25 ou mais recente.

```bash
go run ./cmd/tuitodo
```

## Layout

Quatro painéis:

- **Tarefas** — modos **A fazer** e **Concluídos**
- **Categorias** — inclui a categoria virtual **Todas**
- **Subtarefas** — filhas da tarefa pai selecionada (um nível)
- **Ações** — dicas do painel atual

## Navegação

- Setas escolhem o painel
- **Enter** entra no painel
- **Esc** volta

Os atalhos de teclado não mudam com o idioma. Títulos dos painéis e rótulos seguem o idioma selecionado. Os títulos das tarefas são dados seus.

## Atalhos

- `1` — A fazer
- `2` — Concluídos
- `a` — nova tarefa (em Subtarefas: nova subtarefa)
- `c` — nova categoria
- `s` — Config
- `e` — editar
- `d` — apagar
- `Espaço` — concluir (em Concluídos: reabrir)
- `?` — ajuda de atalhos
- `q` — sair

## Concluídos

A lista da esquerda mostra **só datas** (o dia local da conclusão do pai). O painel da direita lista as tarefas pai daquele dia e as subtarefas. **Espaço** reabre o item em foco. Reabrir uma subtarefa devolve a família para A fazer.

## Idioma

Na primeira execução o app segue o SO (`en-US` ou `pt-BR`, com fallback para `en-US`). Pressione `s` para abrir Config. Trocar o idioma na lista aplica na hora. **Salvar** grava; **Cancelar** ou **Esc** reverte.

## Dados

As tarefas ficam num ficheiro SQLite local:

- macOS: `~/Library/Application Support/tuitodo/tuitodo.db`
- Linux: `~/.local/share/tuitodo/tuitodo.db` (ou `$XDG_DATA_HOME/tuitodo/tuitodo.db`)
```

- [x] **Step 2: Verify Portuguese guide and cross-links**

Run:

```bash
test -f docs/pt-BR/how-to-use.md
grep -F '## Executar' docs/pt-BR/how-to-use.md
grep -F '## Layout' docs/pt-BR/how-to-use.md
grep -F '## Navegação' docs/pt-BR/how-to-use.md
grep -F '## Atalhos' docs/pt-BR/how-to-use.md
grep -F '## Concluídos' docs/pt-BR/how-to-use.md
grep -F '## Idioma' docs/pt-BR/how-to-use.md
grep -F '## Dados' docs/pt-BR/how-to-use.md
grep -F 'go run ./cmd/tuitodo' docs/pt-BR/how-to-use.md
grep -F '[English](../en-US/how-to-use.md)' docs/pt-BR/how-to-use.md
grep -F 'Portuguese (Brazil)' docs/pt-BR/how-to-use.md
! grep -F 'docs/superpowers' docs/pt-BR/how-to-use.md README.md docs/en-US/how-to-use.md
test -f README.md && test -f docs/en-US/how-to-use.md && test -f docs/pt-BR/how-to-use.md
! test -d docs/en-US/how-to-use
! test -d docs/pt-BR/how-to-use
! test -f README.pt-BR.md
```

Expected: all commands succeed (exit 0).

- [x] **Step 3: Commit**

```bash
git add docs/pt-BR/how-to-use.md
git commit -m "$(cat <<'EOF'
docs: add Portuguese (Brazil) how-to-use guide

EOF
)"
```
