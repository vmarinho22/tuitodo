# README e how-to-use

Data: 2026-09-18

Não há README. O utilizador precisa de um ponto de entrada curto e de um guia de uso por idioma. O README não ensina a correr o app; aponta para `docs/<locale>/how-to-use.md`.

## Objetivo

README em inglês com título, contexto e links. O “como usar” vive em guias completos em inglês e português do Brasil.

## Ficheiros

```text
README.md
docs/en-US/how-to-use.md
docs/pt-BR/how-to-use.md
```

`docs/superpowers/` continua só para specs internas. O README não liga para lá.

## README.md (inglês)

Curto. Sem instalar, sem atalhos, sem badges, sem secção de desenvolvimento.

- Título: `tuitodo`
- Um parágrafo: TUI de tarefas (tview), pais com um nível de subtarefas, categorias, SQLite local. Idioma segue o SO (`en-US` / `pt-BR`) e muda na app.
- Secção **How to use** com os dois links:
  - [English](docs/en-US/how-to-use.md)
  - [Portuguese (Brazil)](docs/pt-BR/how-to-use.md)

## How-to-use

Os dois ficheiros têm a **mesma estrutura**, cada um no seu idioma. No topo: troca English · Portuguese (Brazil) para o outro guia.

1. **Run** — Go 1.25+; `go run ./cmd/tuitodo`
2. **Layout** — Tarefas (A fazer / Concluídos), Categorias (Todas), Subtarefas, Ações
3. **Navigation** — setas escolhem o painel; Enter entra; Esc volta
4. **Shortcuts** — `1` / `2`, `a`, `c`, `s`, `e`, `d`, espaço, `?`, `q`
5. **Completed** — esquerda só datas; direita pais + subs; espaço reabre
6. **Language** — deteta o SO; `s` Config; a lista aplica na hora; Salvar grava; Cancelar/Esc reverte
7. **Data** — SQLite local: macOS Application Support (`tuitodo/tuitodo.db`); Linux XDG (`~/.local/share/tuitodo/tuitodo.db` ou `$XDG_DATA_HOME`)

Atalhos de teclado não mudam com o idioma; os rótulos da UI sim. Títulos das tarefas são dados do utilizador.

## Fora de escopo

- README em português ou `README.pt-BR.md`
- Instalar/atalhos no README
- Pasta `how-to-use/` com vários ficheiros
- Documentar specs em `docs/superpowers/`
- Screenshots
