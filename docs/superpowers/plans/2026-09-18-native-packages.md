# Native packages Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking. Do not git add/commit/push (project rule).

**Goal:** Tag `v*` runs GoReleaser v2 to publish macOS/Linux binaries, Ubuntu `.deb`, and a Homebrew formula in `Formula/` of this repo.

**Architecture:** `.goreleaser.yaml` defines builds (`CGO_ENABLED=0`, darwin/linux amd64/arm64), nFPM `.deb`, and `brews` writing `Formula/tuitodo.rb` back to `vmarinho22/tuitodo`. GitHub Actions on tags `v*` runs `goreleaser release --clean` with `GITHUB_TOKEN` only.

**Tech Stack:** GoReleaser v2, nFPM, GitHub Actions (`goreleaser/goreleaser-action@v7`), Homebrew tap via `Formula/`.

## Global Constraints

- No second repo `homebrew-tuitodo`
- No `HOMEBREW_TAP_TOKEN`; only `GITHUB_TOKEN`
- `brew tap vmarinho22/tuitodo https://github.com/vmarinho22/tuitodo`
- `.deb` name `tuitodo_<version>_linux_<arch>.deb`
- No version flag, no LICENSE, no PPA, no Windows
- Do not commit

## File map

| File | Responsibility |
|------|----------------|
| `.goreleaser.yaml` | Builds, nFPM, brews |
| `.github/workflows/release.yml` | Tag-triggered release |
| `.gitignore` | Ignore `dist/` |
| `README.md` | Install section |
| `docs/en-US/how-to-use.md` | Install + Run from source |
| `docs/pt-BR/how-to-use.md` | Instalar + Executar a partir do código |

---

### Task 1: GoReleaser config

**Files:**
- Create: `.goreleaser.yaml`
- Modify: `.gitignore`

- [ ] **Step 1: Add dist/ to .gitignore**

Append:

```
/dist/
```

- [ ] **Step 2: Write .goreleaser.yaml**

```yaml
version: 2

project_name: tuitodo

builds:
  - id: tuitodo
    main: ./cmd/tuitodo
    binary: tuitodo
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w

nfpms:
  - package_name: tuitodo
    file_name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    formats:
      - deb
    bindir: /usr/bin
    section: utils
    description: Local-first terminal todo app
    homepage: https://github.com/vmarinho22/tuitodo
    maintainer: vmarinho22

brews:
  - name: tuitodo
    directory: Formula
    homepage: https://github.com/vmarinho22/tuitodo
    description: Local-first terminal todo app
    repository:
      owner: vmarinho22
      name: tuitodo
      token: "{{ .Env.GITHUB_TOKEN }}"
```

- [ ] **Step 3: Validate config**

```bash
go test ./...
go run github.com/goreleaser/goreleaser/v2@latest check
```

Expected: tests pass; `goreleaser check` exits 0.

---

### Task 2: GitHub Actions release workflow

**Files:**
- Create: `.github/workflows/release.yml`

- [ ] **Step 1: Write the workflow**

```yaml
name: release

on:
  push:
    tags:
      - "v*"

permissions:
  contents: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v5
        with:
          fetch-depth: 0
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.25"
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          distribution: goreleaser
          version: "~> v2"
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

Use `goreleaser-action@v7` if it exists; otherwise `@v6` as documented. Prefer `@v6` with `version: "~> v2"` if v7 is unavailable.

- [ ] **Step 2: Verify workflow YAML**

```bash
grep -F 'tags:' .github/workflows/release.yml
grep -F 'v*' .github/workflows/release.yml
grep -F 'contents: write' .github/workflows/release.yml
grep -F 'GITHUB_TOKEN' .github/workflows/release.yml
! grep -F 'HOMEBREW_TAP_TOKEN' .github/workflows/release.yml
```

Expected: greps succeed; last command exits 0.

---

### Task 3: Install docs

**Files:**
- Modify: `README.md`
- Modify: `docs/en-US/how-to-use.md`
- Modify: `docs/pt-BR/how-to-use.md`

- [ ] **Step 1: README Install section** after the language paragraph, before How to use:

```markdown
## Install

macOS (Homebrew):

```bash
brew tap vmarinho22/tuitodo https://github.com/vmarinho22/tuitodo
brew install tuitodo
```

Ubuntu (from the GitHub Release `.deb`):

```bash
sudo apt install ./tuitodo_<version>_linux_amd64.deb
```
```

- [ ] **Step 2: English how-to-use** — after the language switcher, add **Install** (Homebrew + `.deb`); rename **Run** to **Run from source** (keep Go 1.25 and `go run ./cmd/tuitodo`).

- [ ] **Step 3: Portuguese how-to-use** — same: **Instalar**, then **Executar a partir do código**.

- [ ] **Step 4: Verify**

```bash
grep -F 'brew tap vmarinho22/tuitodo https://github.com/vmarinho22/tuitodo' README.md docs/en-US/how-to-use.md docs/pt-BR/how-to-use.md
grep -F '## Run from source' docs/en-US/how-to-use.md
grep -F '## Executar a partir do código' docs/pt-BR/how-to-use.md
! grep -F 'homebrew-tuitodo' README.md docs/en-US/how-to-use.md docs/pt-BR/how-to-use.md .goreleaser.yaml
```
