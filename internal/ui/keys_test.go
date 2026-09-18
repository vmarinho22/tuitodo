package ui

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestIsEscapeKey(t *testing.T) {
	cases := []struct {
		name  string
		event *tcell.EventKey
		want  bool
	}{
		{
			name:  "escape key",
			event: tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone),
			want:  true,
		},
		{
			name:  "escape rune",
			event: tcell.NewEventKey(tcell.KeyRune, 27, tcell.ModNone),
			want:  true,
		},
		{
			name:  "enter",
			event: tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone),
			want:  false,
		},
		{
			name:  "letter",
			event: tcell.NewEventKey(tcell.KeyRune, 'q', tcell.ModNone),
			want:  false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got := isEscapeKey(testCase.event)
			if got != testCase.want {
				t.Fatalf("isEscapeKey() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestShortcutRuneIgnoresCapsLock(t *testing.T) {
	got := shortcutRune(tcell.NewEventKey(tcell.KeyRune, 'Q', tcell.ModNone))
	if got != 'q' {
		t.Fatalf("shortcutRune('Q') = %q, want 'q'", got)
	}
	got = shortcutRune(tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone))
	if got != 'a' {
		t.Fatalf("shortcutRune('a') = %q, want 'a'", got)
	}
}

func TestHelpBodyHeightFitsSmallTerminals(t *testing.T) {
	if got := helpBodyHeight(1); got != 2 {
		t.Fatalf("min height = %d, want 2", got)
	}
	if got := helpBodyHeight(8); got != 8 {
		t.Fatalf("short body = %d, want 8", got)
	}
	if got := helpBodyHeight(30); got != helpBodyMaxHeight {
		t.Fatalf("long body = %d, want cap %d", got, helpBodyMaxHeight)
	}
}

func TestColorHelpTitlesHighlightsSectionHeaders(t *testing.T) {
	got := colorHelpTitles("Navegação\n  [a]         nova tarefa\nModos")
	if !strings.Contains(got, "[yellow::b]") || !strings.Contains(got, "Navegação") {
		t.Fatalf("title not highlighted: %q", got)
	}
	if strings.Contains(got, "[yellow::b]  [") {
		t.Fatal("shortcut lines should not use the title color")
	}
	if !strings.Contains(got, tview.Escape("  [a]         nova tarefa")) {
		t.Fatalf("shortcut line must be escaped for tview: %q", got)
	}
}
