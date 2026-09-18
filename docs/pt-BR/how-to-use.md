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

As tarefas ficam em um arquivo SQLite local:

- macOS: `~/Library/Application Support/tuitodo/tuitodo.db`
- Linux: `~/.local/share/tuitodo/tuitodo.db` (ou `$XDG_DATA_HOME/tuitodo/tuitodo.db`)
