package ui

import "github.com/gdamore/tcell/v2"

func isEscapeKey(event *tcell.EventKey) bool {
	if event.Key() == tcell.KeyEscape {
		return true
	}
	return event.Key() == tcell.KeyRune && event.Rune() == 27
}
