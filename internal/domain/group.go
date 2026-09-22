package domain

import (
	"sort"
	"time"
)

type CompletedDayEntry struct {
	Parent   Task
	Subtasks []Task
}

type CompletedDay struct {
	Date    time.Time
	Entries []CompletedDayEntry
}

type dayKey struct {
	year  int
	month time.Month
	day   int
}

func GroupCompletedActivityByDay(parents []Task, subtasksByParent map[int64][]Task, location *time.Location) []CompletedDay {
	type entryBuilder struct {
		parent   Task
		subtasks []Task
	}
	byDay := make(map[dayKey]map[int64]*entryBuilder)

	ensure := func(key dayKey, parent Task) *entryBuilder {
		if byDay[key] == nil {
			byDay[key] = make(map[int64]*entryBuilder)
		}
		if byDay[key][parent.ID] == nil {
			byDay[key][parent.ID] = &entryBuilder{parent: parent}
		}
		return byDay[key][parent.ID]
	}

	for _, parent := range parents {
		if parent.CompletedAt != nil {
			ensure(calendarDayKey(*parent.CompletedAt, location), parent)
		}
		for _, subtask := range subtasksByParent[parent.ID] {
			if subtask.CompletedAt == nil {
				continue
			}
			entry := ensure(calendarDayKey(*subtask.CompletedAt, location), parent)
			entry.subtasks = append(entry.subtasks, subtask)
		}
	}

	keys := make([]dayKey, 0, len(byDay))
	for key := range byDay {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		left := time.Date(keys[i].year, keys[i].month, keys[i].day, 0, 0, 0, 0, location)
		right := time.Date(keys[j].year, keys[j].month, keys[j].day, 0, 0, 0, 0, location)
		return left.After(right)
	})

	days := make([]CompletedDay, 0, len(keys))
	for _, key := range keys {
		entries := make([]CompletedDayEntry, 0)
		for _, parent := range parents {
			builder := byDay[key][parent.ID]
			if builder == nil {
				continue
			}
			entries = append(entries, CompletedDayEntry{
				Parent:   builder.parent,
				Subtasks: builder.subtasks,
			})
		}
		days = append(days, CompletedDay{
			Date:    time.Date(key.year, key.month, key.day, 0, 0, 0, 0, location),
			Entries: entries,
		})
	}
	return days
}

func calendarDayKey(moment time.Time, location *time.Location) dayKey {
	local := moment.In(location)
	return dayKey{year: local.Year(), month: local.Month(), day: local.Day()}
}
