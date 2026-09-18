package ui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNextPane(t *testing.T) {
	cases := []struct {
		name    string
		current paneID
		lastTop paneID
		key     tcell.Key
		want    paneID
	}{
		{name: "tasks down", current: paneTasks, lastTop: paneTasks, key: tcell.KeyDown, want: paneCategories},
		{name: "tasks right", current: paneTasks, lastTop: paneTasks, key: tcell.KeyRight, want: paneDetail},
		{name: "tasks up", current: paneTasks, lastTop: paneTasks, key: tcell.KeyUp, want: paneActions},
		{name: "categories up", current: paneCategories, lastTop: paneCategories, key: tcell.KeyUp, want: paneTasks},
		{name: "categories down", current: paneCategories, lastTop: paneCategories, key: tcell.KeyDown, want: paneActions},
		{name: "categories right", current: paneCategories, lastTop: paneCategories, key: tcell.KeyRight, want: paneDetail},
		{name: "detail left remembers categories", current: paneDetail, lastTop: paneCategories, key: tcell.KeyLeft, want: paneCategories},
		{name: "detail left defaults to tasks", current: paneDetail, lastTop: paneTasks, key: tcell.KeyLeft, want: paneTasks},
		{name: "detail down", current: paneDetail, lastTop: paneDetail, key: tcell.KeyDown, want: paneActions},
		{name: "actions up to last top", current: paneActions, lastTop: paneDetail, key: tcell.KeyUp, want: paneDetail},
		{name: "actions down", current: paneActions, lastTop: paneDetail, key: tcell.KeyDown, want: paneTasks},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got := nextPane(testCase.current, testCase.lastTop, testCase.key)
			if got != testCase.want {
				t.Fatalf("nextPane() = %v, want %v", got, testCase.want)
			}
		})
	}
}
