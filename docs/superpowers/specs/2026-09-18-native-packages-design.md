# Pacotes nativos: Homebrew tap e .deb

Data: 2026-09-18

O tuitodo é um único binário Go sem CGO (`modernc.org/sqlite`). Não há CI nem empacotamento. O utilizador instala via Homebrew (macOS) e `.deb` nas GitHub Releases (Ubuntu). Sem PPA.

## Objetivo

Num tag `v*`, o GitHub Actions corre GoReleaser v2: binários, `.deb` e cask Homebrew em `Casks/tuitodo.rb` **neste** repositório.

## Repositório

Só `vmarinho22/tuitodo`. Não há `homebrew-tuitodo`.

O tap é a pasta `Casks/` na raiz (GoReleaser v2 gera casks, não fórmulas). `brew tap` clona o Git inteiro; não dá para tapar só um subdirectório.

Instalação:

```bash
brew tap vmarinho22/tuitodo https://github.com/vmarinho22/tuitodo
brew install --cask tuitodo
```

A URL extra é obrigatória: sem ela o Homebrew procura `vmarinho22/homebrew-tuitodo`.

## Fluxo

1. `git tag v0.1.0 && git push origin v0.1.0`
2. Workflow `release` (só tags `v*`)
3. GoReleaser gera artefactos, cria a GitHub Release, faz commit de `Casks/tuitodo.rb` na branch por omissão

`GITHUB_TOKEN` (com `contents: write`) chega para a Release e para o commit da fórmula. Sem PAT extra.

## Builds

`.goreleaser.yaml` (GoReleaser v2):

- `main: ./cmd/tuitodo`, `binary: tuitodo`
- `CGO_ENABLED=0`
- `goos: [linux, darwin]`
- `goarch: [amd64, arm64]`
- Sem ldflags de versão (a app não tem `-v`)

## .deb (nFPM)

- `formats: [deb]`
- `bindir: /usr/bin`
- `section: utils`
- Nome: `tuitodo_<version>_linux_<arch>.deb`
- Sem `.rpm`, `.apk`, Snap, Flatpak

Ubuntu: descarregar o `.deb` da Release e `sudo apt install ./tuitodo_<version>_linux_amd64.deb` (ou `arm64`).

## Homebrew

Bloco `homebrew_casks` com `repository.owner: vmarinho22`, `name: tuitodo`, `directory: Casks`, `token: "{{ .Env.GITHUB_TOKEN }}"`.

Homepage: `https://github.com/vmarinho22/tuitodo`. Descrição alinhada ao README (local-first terminal todo). Sem `license` na fórmula enquanto o repo não tiver `LICENSE` (este trabalho não adiciona licença).

## CI

`.github/workflows/release.yml`:

- `on.push.tags: ["v*"]`
- `permissions.contents: write`
- `actions/checkout` com `fetch-depth: 0`
- `actions/setup-go` com Go 1.25
- `goreleaser/goreleaser-action@v7`, `version: "~> v2"`, `args: release --clean`
- `GITHUB_TOKEN` apenas

Sem snapshot em pull request. Sem secret `HOMEBREW_TAP_TOKEN`.

## Docs

README (inglês), secção **Install** após o resumo, antes de **How to use**:

```bash
brew tap vmarinho22/tuitodo https://github.com/vmarinho22/tuitodo
brew install --cask tuitodo
```

```bash
sudo apt install ./tuitodo_<version>_linux_amd64.deb
```

`docs/en-US/how-to-use.md` e `docs/pt-BR/how-to-use.md`: nova secção **Install** no topo (depois do language switcher); **Run** passa a **Run from source** / **Executar a partir do código** com `go run ./cmd/tuitodo`. Sem badges.

## Fora de escopo

- Segundo repo `homebrew-tuitodo`
- PPA / Launchpad
- `.pkg`, notarização, Homebrew Core
- Windows
- Flag de versão no binário
- Ficheiro `LICENSE`
