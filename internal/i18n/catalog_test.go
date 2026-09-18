package i18n

import (
	"errors"
	"strings"
	"testing"
	"time"

	"tuitodo/internal/domain"
	"tuitodo/internal/store"
)

func TestTUsesCurrentLocale(t *testing.T) {
	catalog := NewCatalog(LocalePtBR)
	got := catalog.T(KeyTasksTitle)
	if got != "Tarefas" {
		t.Fatalf("T(tasks.title) = %q, want Tarefas", got)
	}
}

func TestTFallsBackToEnglish(t *testing.T) {
	catalog := NewCatalog(LocalePtBR)
	catalog.messages = map[Locale]map[string]string{
		LocalePtBR: {},
		LocaleEnUS: catalogs[LocaleEnUS],
	}
	got := catalog.T(KeyTasksTitle)
	if got != "Tasks" {
		t.Fatalf("fallback = %q, want Tasks", got)
	}
}

func TestTMissingKeyReturnsKey(t *testing.T) {
	catalog := NewCatalog(LocaleEnUS)
	got := catalog.T("missing.key")
	if got != "missing.key" {
		t.Fatalf("T(missing) = %q", got)
	}
}

func TestTFormatsArgs(t *testing.T) {
	catalog := NewCatalog(LocaleEnUS)
	got := catalog.T(KeyDeleteTaskConfirm, "Gym")
	if got != `Delete "Gym"?` {
		t.Fatalf("formatted = %q", got)
	}
}

func TestShortcutStringsPutKeysInBrackets(t *testing.T) {
	en := NewCatalog(LocaleEnUS)
	if got := en.T(KeyShortcutsPending); got != " [a] add  [e] edit  [d] delete  [space] complete" {
		t.Fatalf("en pending = %q", got)
	}
	if got := en.T(KeyShortcutsDone); got != " [space] reopen" {
		t.Fatalf("en done = %q", got)
	}
	pt := NewCatalog(LocalePtBR)
	if got := pt.T(KeyShortcutsPending); got != " [a] adicionar  [e] editar  [d] apagar  [espaço] concluir" {
		t.Fatalf("pt pending = %q", got)
	}
}

func TestHelpBodyIsGroupedBySection(t *testing.T) {
	pt := NewCatalog(LocalePtBR)
	body := pt.T(KeyHelpBody)
	for _, want := range []string{"Navegação", "Modos", "Geral", "Tarefas", "Concluídos", "[espaço]"} {
		if !strings.Contains(body, want) {
			t.Fatalf("help body missing %q:\n%s", want, body)
		}
	}
	if strings.Contains(body, "Setas escolhem o painel  Enter") {
		t.Fatal("help body still uses the packed one-line layout")
	}
	en := NewCatalog(LocaleEnUS)
	if !strings.Contains(en.T(KeyHelpBody), "Navigation") {
		t.Fatalf("en help missing Navigation:\n%s", en.T(KeyHelpBody))
	}
}

func TestErrorKnownSentinels(t *testing.T) {
	catalog := NewCatalog(LocalePtBR)
	if got := catalog.Error(domain.ErrEmptyTitle); got != "O título não pode ser vazio." {
		t.Fatalf("empty title = %q", got)
	}
	if got := catalog.Error(store.ErrCategoryInUse); got != "Essa categoria ainda tem tarefas." {
		t.Fatalf("category in use = %q", got)
	}
	catalog.SetLocale(LocaleEnUS)
	if got := catalog.Error(domain.ErrEmptyTitle); got != "Title cannot be empty." {
		t.Fatalf("empty title en = %q", got)
	}
}

func TestErrorUnknownUsesGeneric(t *testing.T) {
	catalog := NewCatalog(LocaleEnUS)
	got := catalog.Error(errors.New("sqlite boom"))
	if got != "Something went wrong." {
		t.Fatalf("unknown = %q", got)
	}
}

func TestFormatDateAndTime(t *testing.T) {
	moment := time.Date(2026, 9, 18, 8, 25, 0, 0, time.Local)
	en := NewCatalog(LocaleEnUS)
	if got := en.FormatDate(moment); got != "2026-09-18" {
		t.Fatalf("en date = %q", got)
	}
	if got := en.FormatTime(moment); got != "8:25 AM" {
		t.Fatalf("en time = %q", got)
	}
	pt := NewCatalog(LocalePtBR)
	if got := pt.FormatDate(moment); got != "18/09/2026" {
		t.Fatalf("pt date = %q", got)
	}
	if got := pt.FormatTime(moment); got != "08:25" {
		t.Fatalf("pt time = %q", got)
	}
}
