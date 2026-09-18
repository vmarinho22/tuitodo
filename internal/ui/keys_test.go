package ui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
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
