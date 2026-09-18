# Concluídos por data

Data: 2026-09-18

Em Concluídos, misturar data e tarefa na mesma lista deixa o dia pouco claro. A lista esquerda passa a mostrar só datas. O painel Subtarefas mostra as tarefas pai daquele dia, cada uma com suas subtarefas embaixo. Pendentes não muda. Regras de domínio (cascata, um nível, agrupamento por dia local) não mudam.

## Objetivo

Consultar o que foi concluído por dia, sem confundir cabeçalho de data com item de tarefa. Dá para reabrir; não dá para criar, editar ou apagar nesse modo.

## Pendentes

Igual ao spec [2026-09-18-subtarefas-panel-design.md](2026-09-18-subtarefas-panel-design.md): Tarefas lista pais; Subtarefas mostra o título do pai, a lista de subs e o rodapé `a adicionar  e editar  d apagar  espaço concluir`.

## Layout em Concluídos (`2`)

```text
┌─ Tarefas · Concluídos ──┐  ┌─ Subtarefas ──────────────┐
│ 18/09/2026              │  │ 18/09/2026                │
│ 17/09/2026              │  │                           │
│                         │  │ 08:25  CE-9915 Reemitir…  │
│                         │  │   [x] anexar receita      │
│                         │  │   [x] enviar ao paciente  │
│                         │  │ 09:10  Academia           │
│                         │  │                           │
│                         │  │ espaço reabrir            │
└─────────────────────────┘  └───────────────────────────┘
```

### Esquerda (Tarefas)

- Só datas, formato `02/01/2006`, mais recente primeiro.
- Sem horário, sem título, sem categoria.
- Datas vêm de `GroupParentTasksByCompletedDay` (dia local de `completed_at` do pai).
- Filtro de categoria continua: só entram dias que têm pelo menos um pai da categoria (ou Todas).

### Direita (Subtarefas)

- Cabeçalho: a data selecionada. Sem pai selecionado / lista vazia: `Selecione um dia`.
- Lista: para cada pai do dia, na ordem do grupo (já é a ordem de conclusão no store):
  - Linha do pai: `15:04  título` (sem checkbox).
  - Subs indentadas: `  [x] título` (todas concluídas neste modo, salvo inconsistência).
  - Pai sem sub: só a linha do pai.
- Rodapé: `espaço reabrir` apenas. Sem `a` / `e` / `d`.

## Navegação

Setas escolhem o painel; Enter entra; Esc volta. Clique também entra.

Depois de Enter em Tarefas: setas / `j` `k` mudam o dia; a direita atualiza na hora.

Depois de Enter em Subtarefas: setas / `j` `k` percorrem linhas de pai e de sub.

`1` / `2` trocam Pendentes / Concluídos de qualquer lugar. Ao voltar a Concluídos, restaura o último dia se ainda existir; senão o mais recente.

## Teclado em Concluídos (painel ativo)

### Tarefas (lista de datas)

`a` / `e` / `d` / espaço não fazem nada nesta lista.

### Subtarefas (conteúdo do dia)

| Foco | Espaço |
| --- | --- |
| Pai **sem** subs | Reabre o pai (`ReopenParentWithoutSubtasks`); família some de Concluídos e volta a Pendentes. |
| Pai **com** subs | Não faz nada (regra atual: reabrir o pai com subs auto-concluiria de novo). |
| Sub | `ReopenSubtask`: aquela sub e o pai ficam pendentes; outras subs concluídas mantêm `completed_at`; a família volta a Pendentes. |

`a` / `e` / `d` não fazem nada neste modo.

## Implementação

- Em `showingCompleted`, `taskListEntries` só tem datas (`isDayHeader`); `selectedDay` (data local) substitui a seleção de pai como âncora da direita.
- `refreshDetail` ramifica: Pendentes = spec de subtarefas; Concluídos = cabeçalho da data + linhas pai/sub + rodapé `espaço reabrir`.
- `toggleFocusedCompletion` no modo concluído só trata reabertura conforme a linha focada; ignora `a`/`e`/`d`.
- `GroupParentTasksByCompletedDay` permanece no domínio.
- Sem mudança de schema SQLite.

## Testes

- Testes atuais de agrupamento e de reabrir continuam passando.
- Opcional: teste de UI não é obrigatório. Se houver helper puro que monta as linhas do dia (pai + subs), testar ordem e indentação.

## Fora de escopo

- Editar, apagar ou adicionar em Concluídos.
- Mostrar categoria na linha do pai concluído (pode ir no texto secundário depois, se fizer falta).
- Calendário gráfico.
- Mudar Pendentes.
