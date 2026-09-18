# Abas Pendentes / Concluídos no painel Tarefas

Data: 2026-09-18

O título `Tarefas · Pendentes` não mostra que existe outro modo nem como chegar lá. O painel Tarefas passa a ter o nome fixo na borda e, logo abaixo, as duas opções com atalho. Não há widget de abas no tview; a linha é só leitura. `1` / `2` continuam globais.

## Objetivo

Deixar óbvio em qual modo estamos e que dá para ir ao outro, sem criar um controle novo nem mudar a navegação de painéis.

## Layout

O painel esquerdo de cima vira um `Flex` vertical com borda e título `Tarefas` (igual a Subtarefas):

```text
┌─ Tarefas ─────────────────┐
│ Pendentes 1  Concluídos 2 │
│                           │
│ CE-9915 Reemitir…         │
│ Academia                  │
└───────────────────────────┘
```

1. **Borda** — título fixo `Tarefas`. A cor da borda segue o painel selecionado / ativo, como hoje.
2. **Linha de modos** — `TextView` de 1 linha, só leitura, sem foco, sem clique para trocar. Ativa em negrito (`[::b]`), inativa apagada (`[::d]`).
3. **Lista** — pendentes ou datas de concluídos, como hoje. Único filho focável.

Pendentes: ` [::b]Pendentes 1[::-]  [::d]Concluídos 2[::-] `  
Concluídos: o inverso.

## Navegação e teclado

Setas escolhem o painel; Enter entra na lista; Esc volta. Clique no painel (linha de modos ou lista) entra em Tarefas e foca a lista.

`1` e `2` continuam globais, de qualquer painel. A linha só reflete o modo; não recebe `1`/`2` sozinha e não é clicável para trocar.

## Implementação

- `tasksPane` (`Flex`) envolve a linha de modos e `taskList`.
- `panePrimitive` / `paneBox` de Tarefas apontam para `tasksPane`.
- Ao entrar no painel, o foco vai para `taskList` (como Subtarefas → lista de subs).
- Função pura monta o texto da linha nos dois modos.
- Sem mudança de domínio, store ou atalhos.

## Testes

- Texto da linha em Pendentes e em Concluídos (tags e atalhos).

## Fora de escopo

- Abas clicáveis.
- `Pages` / `Button` para os modos.
- Tirar `1` / `2` do teclado global.
- Mudar o conteúdo de Pendentes ou Concluídos.
