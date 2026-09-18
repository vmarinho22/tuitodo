package ui

import (
	"testing"
	"time"

	"tuitodo/internal/domain"
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

func TestBuildCompletedDayLinesOrdersParentThenSubs(t *testing.T) {
	completedAt := time.Date(2026, 9, 18, 8, 25, 0, 0, time.Local)
	parentWithSubs := domain.Task{ID: 1, Title: "CE-9915 Reemitir", CompletedAt: &completedAt}
	parentAlone := domain.Task{ID: 2, Title: "Academia", CompletedAt: ptrTime(time.Date(2026, 9, 18, 9, 10, 0, 0, time.Local))}
	subtask := domain.Task{ID: 11, ParentID: ptrInt64(1), Title: "anexar receita", CompletedAt: &completedAt}

	lines := buildCompletedDayLines(domain.CompletedDay{
		Date:  time.Date(2026, 9, 18, 0, 0, 0, 0, time.Local),
		Tasks: []domain.Task{parentWithSubs, parentAlone},
	}, map[int64][]domain.Task{
		1: {subtask},
	})

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

func ptrTime(value time.Time) *time.Time {
	return &value
}

func ptrInt64(value int64) *int64 {
	return &value
}
