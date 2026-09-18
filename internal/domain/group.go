package domain

import "time"

type CompletedDay struct {
	Date  time.Time
	Tasks []Task
}

func GroupParentTasksByCompletedDay(parentTasks []Task, location *time.Location) []CompletedDay {
	type dayKey struct {
		year  int
		month time.Month
		day   int
	}

	groups := make(map[dayKey][]Task)
	order := make([]dayKey, 0)

	for _, parentTask := range parentTasks {
		if parentTask.CompletedAt == nil {
			continue
		}
		localCompletedAt := parentTask.CompletedAt.In(location)
		key := dayKey{
			year:  localCompletedAt.Year(),
			month: localCompletedAt.Month(),
			day:   localCompletedAt.Day(),
		}
		if _, exists := groups[key]; !exists {
			order = append(order, key)
		}
		groups[key] = append(groups[key], parentTask)
	}

	days := make([]CompletedDay, 0, len(order))
	for _, key := range order {
		days = append(days, CompletedDay{
			Date:  time.Date(key.year, key.month, key.day, 0, 0, 0, 0, location),
			Tasks: groups[key],
		})
	}
	return days
}
