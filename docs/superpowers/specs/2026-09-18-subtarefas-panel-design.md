# Painel Subtarefas

Data: 2026-09-18

O painel direito deixa de se chamar Detalhe e passa a ser o lugar de gerenciar subtarefas da tarefa pai selecionada em Tarefas. A navegação por seções (setas escolhem, Enter entra, Esc volta) permanece. Domínio, SQLite e regras de conclusão em cascata não mudam.

## Objetivo

Criar e gerenciar subtarefas no painel direito, com atalhos visíveis só quando o painel está ativo.

## Layout

O painel direito é um `Flex` vertical com borda e título `Subtarefas`:

```text
┌─ Subtarefas ──────────────────────────┐
│ Relatório Q3                          │
│                                       │
│ [ ] escrever intro                    │
│ [x] revisar números                   │
│                                       │
│ a adicionar  e editar  d apagar  espaço concluir │
└───────────────────────────────────────┘
```

1. **Cabeçalho** — título da tarefa pai, só leitura. Não recebe foco.
2. **Lista** — só subtarefas, checkbox `[ ]` / `[x]` como hoje. É o único filho focável do painel.
3. **Rodapé** — atalhos fixos: `a adicionar  e editar  d apagar  espaço concluir`. Não recebe foco.

Título da borda não inclui mais categoria · pai. A categoria continua na lista de Tarefas (texto secundário).

### Estados vazios

- Nenhum pai selecionado em Tarefas: cabeçalho `Selecione uma tarefa`, lista vazia, rodapé visível.
- Pai sem subtarefas: cabeçalho com o título do pai, lista vazia, rodapé visível (`a` cria a primeira).

## Navegação

Setas (e `h` `j` `k` `l`) escolhem a seção. Enter entra. Esc volta a escolher seção. Clique no painel também entra.

A seção direita continua uma só unidade para as setas. Depois de Enter, o foco vai para a lista de subtarefas; setas / `j` `k` movem o item.

## Teclado depois de Enter

### Painel Subtarefas

| Tecla | Ação |
| --- | --- |
| `a` | Modal Nova subtarefa (título; sem categoria; herda o pai) |
| `e` | Modal editar título da sub selecionada |
| `d` | Confirmar apagar a sub selecionada |
| espaço | Concluir ou reabrir a sub selecionada |

Reabrir uma sub ainda devolve a família a Pendentes; as outras subs mantêm `completed_at`. Pai já concluído: `a` mostra o erro de reabrir antes de adicionar. Sem item selecionado: `e` / `d` / espaço não fazem nada; `a` ainda cria.

Fora deste painel (incluindo modo “escolher seção”), `a` / `e` / `d` / espaço **não** gerenciam subtarefa.

### Painel Tarefas (depois de Enter)

| Tecla | Ação |
| --- | --- |
| `e` | Editar título do pai |
| `d` | Confirmar apagar o pai (e as subs, cascade) |
| espaço | Concluir ou reabrir o pai (mesmas regras atuais: pai com subs concluído não reabre pelo pai) |
| `a` | Não cria subtarefa (abre nova tarefa pai, atalho global) |

### Outros

- **Categorias** depois de Enter: `e` renomeia, `d` apaga, como hoje.
- **Ações** (barra inferior): continua global — `a` tarefa, `c` categoria, `?`, `q`, etc.
- Modais: Esc fecha o modal; o card compacto de apagar/erro e o card grande de formulário não mudam.

## Regras que não mudam

- Um nível só: sub não tem sub.
- Sub herda categoria do pai (não grava `category_id` próprio).
- Concluir o pai conclui as subs pendentes; concluir a última sub conclui o pai.
- Concluídos agrupam pelo dia local de `completed_at` do pai.

## Implementação

- Trocar o `detailList` solto por um container (`subtasksPane`) com cabeçalho, lista e rodapé.
- `paneDetail` aponta para esse container; o foco interno após Enter é a lista.
- `refreshDetail` passa a preencher cabeçalho + lista só de subs.
- `a` / `e` / `d` / espaço em `handleKeys` só gerenciam subtarefa se `paneActive` e o painel for Subtarefas; `a` fora desse painel continua criando tarefa pai.
- Modais de nova sub, editar e apagar reutilizam os cards atuais.
- Sem mudança em `internal/domain` nem `internal/store`, salvo se algum helper de UI precisar filtrar o pai da lista.

## Testes

- Testes de domínio e store existentes devem continuar passando.
- Não é obrigatório teste de TUI. Se houver teste de `nextPane`, ele continua válido com o painel direito como uma seção.

## Fora de escopo

- Subtarefas aninhadas.
- Categoria na sub.
- Inline add sem modal.
- Mudar a barra Ações global.
