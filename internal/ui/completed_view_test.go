package ui

import (
	"testing"
	"time"

	"tuitodo/internal/domain"
	"tuitodo/internal/i18n"

	"github.com/rivo/tview"
)

func TestSameCalendarDay(t *testing.T) {
	morning := time.Date(2026, 9, 18, 8, 25, 0, 0, time.Local)
	evening := time.Date(2026, 9, 18, 21, 10, 0, 0, time.Local)
	nextDay := time.Date(2026, 9, 19, 0, 0, 0, 0, time.Local)

	if !sameCalendarDay(morning, evening) {
		t.Fatal("same calendar day should match across hours")
	}
	if sameCalendarDay(morning, nextDay) {
		t.Fatal("different calendar days should not match")
	}
}

func TestIndexOfSelectedDayRestoresMatchingDate(t *testing.T) {
	selected := time.Date(2026, 9, 17, 0, 0, 0, 0, time.Local)
	app := &App{
		selectedDay: selected,
		taskListEntries: []taskListEntry{
			{isDayHeader: true, day: time.Date(2026, 9, 18, 0, 0, 0, 0, time.Local)},
			{isDayHeader: true, day: time.Date(2026, 9, 17, 8, 0, 0, 0, time.Local)},
		},
	}

	if got := app.indexOfSelectedDay(); got != 1 {
		t.Fatalf("indexOfSelectedDay() = %d, want 1", got)
	}
}

func TestParentDetailTitleIncludesCategory(t *testing.T) {
	got := parentDetailTitle("work", "go to mkt and get some tasks")
	want := "[work] go to mkt and get some tasks"
	if got != want {
		t.Fatalf("parentDetailTitle() = %q, want %q", got, want)
	}
}

func TestParentDetailTitleWithoutCategory(t *testing.T) {
	got := parentDetailTitle("", "go to mkt and get some tasks")
	want := "go to mkt and get some tasks"
	if got != want {
		t.Fatalf("parentDetailTitle() = %q, want %q", got, want)
	}
}

func TestCheckboxLabelPendingEscapesTviewStyleTags(t *testing.T) {
	got := checkboxLabel(domain.Task{Title: "buy milk"})
	want := tview.Escape("[ ] buy milk")
	if got != want {
		t.Fatalf("checkboxLabel() = %q, want %q", got, want)
	}
}

func TestBuildCompletedDayLinesOrdersParentThenSubs(t *testing.T) {
	completedAt := time.Date(2026, 9, 18, 8, 25, 0, 0, time.Local)
	day := time.Date(2026, 9, 18, 0, 0, 0, 0, time.Local)
	parentWithSubs := domain.Task{ID: 1, Title: "CE-9915 Reemitir", CompletedAt: &completedAt}
	parentAlone := domain.Task{ID: 2, Title: "Academia", CompletedAt: ptrTime(time.Date(2026, 9, 18, 9, 10, 0, 0, time.Local))}
	subtask := domain.Task{ID: 11, ParentID: ptrInt64(1), Title: "anexar receita", CompletedAt: &completedAt}

	formatTime := i18n.NewCatalog(i18n.LocalePtBR).FormatTime
	lines := buildCompletedDayLines(domain.CompletedDay{
		Date: day,
		Entries: []domain.CompletedDayEntry{
			{Parent: parentWithSubs, Subtasks: []domain.Task{subtask}},
			{Parent: parentAlone},
		},
	}, formatTime)

	if len(lines) != 3 {
		t.Fatalf("len(lines) = %d, want 3", len(lines))
	}
	if !lines[0].isParent || lines[0].label != "08:25  CE-9915 Reemitir" {
		t.Fatalf("parent line = %+v", lines[0])
	}
	if lines[1].isParent || lines[1].label != "  [✓] anexar receita" {
		t.Fatalf("sub line = %+v", lines[1])
	}
	if !lines[2].isParent || lines[2].label != "09:10  Academia" {
		t.Fatalf("parent without subs = %+v", lines[2])
	}
}

func TestCompletedParentLabelOmitsTimeWhenParentPendingOnDay(t *testing.T) {
	day := time.Date(2026, 9, 18, 0, 0, 0, 0, time.Local)
	formatTime := i18n.NewCatalog(i18n.LocalePtBR).FormatTime
	pending := domain.Task{ID: 1, Title: "pai"}
	if got := completedParentLabel(day, pending, formatTime); got != "pai" {
		t.Fatalf("pending parent = %q", got)
	}
	otherDay := time.Date(2026, 9, 17, 8, 0, 0, 0, time.Local)
	completedOtherDay := domain.Task{ID: 2, Title: "pai", CompletedAt: &otherDay}
	if got := completedParentLabel(day, completedOtherDay, formatTime); got != "pai" {
		t.Fatalf("other-day parent = %q", got)
	}
}

func ptrTime(value time.Time) *time.Time {
	return &value
}

func ptrInt64(value int64) *int64 {
	return &value
}
