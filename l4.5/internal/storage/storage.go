package storage

import (
	"fmt"
	"sync"
	"time"

	"l2.18/internal/model"
)

type Storage struct {
	mu     sync.RWMutex
	events map[int]map[int64]map[string]model.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[int]map[int64]map[string]model.Event),
	}
}

func dayKey(t time.Time) int64 {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location()).Unix()
}

func (st *Storage) Create(userID int, event model.Event) {
	st.mu.Lock()
	defer st.mu.Unlock()

	day := dayKey(event.Date)

	if st.events[userID] == nil {
		st.events[userID] = make(map[int64]map[string]model.Event)
	}
	if st.events[userID][day] == nil {
		st.events[userID][day] = make(map[string]model.Event)
	}

	st.events[userID][day][event.Name] = event
}

func (st *Storage) ExactEventExists(userID int, date time.Time, name string) bool {
	st.mu.RLock()
	defer st.mu.RUnlock()

	day := dayKey(date)

	user := st.events[userID]
	if user == nil {
		return false
	}
	events := user[day]
	if events == nil {
		return false
	}
	_, ok := events[name]
	return ok
}

func (st *Storage) EventAtTimeExists(userID int, date time.Time) bool {
	st.mu.RLock()
	defer st.mu.RUnlock()

	day := dayKey(date)
	user := st.events[userID]
	if user == nil {
		return false
	}
	return len(user[day]) > 0
}

func (st *Storage) Update(userID int, date time.Time, updated model.Event) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	day := dayKey(date)
	user := st.events[userID]
	if user == nil {
		return fmt.Errorf("user not found")
	}
	events := user[day]
	if events == nil {
		return fmt.Errorf("event not found")
	}

	if _, ok := events[updated.Name]; !ok {
		return fmt.Errorf("event not found")
	}

	events[updated.Name] = updated
	return nil
}

func (st *Storage) Delete(userID int, date time.Time, name string) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	day := dayKey(date)
	user := st.events[userID]
	if user == nil {
		return fmt.Errorf("user not found")
	}
	events := user[day]
	if events == nil {
		return fmt.Errorf("event not found")
	}

	if _, ok := events[name]; !ok {
		return fmt.Errorf("event not found")
	}

	delete(events, name)
	return nil
}

func (st *Storage) EventsForDay(userID int, date time.Time) ([]model.Event, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	day := dayKey(date)
	user := st.events[userID]
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	events := user[day]
	if len(events) == 0 {
		return nil, fmt.Errorf("no events for this day")
	}

	result := make([]model.Event, 0, len(events))
	for _, e := range events {
		result = append(result, e)
	}
	return result, nil
}

func (st *Storage) EventsForWeek(userID int, date time.Time) ([]model.Event, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	user := st.events[userID]
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	targetYear, targetWeek := date.ISOWeek()
	result := []model.Event{}

	for _, eventsByName := range user {
		for _, e := range eventsByName {
			year, week := e.Date.ISOWeek()
			if year == targetYear && week == targetWeek {
				result = append(result, e)
			}
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no events for this week")
	}
	return result, nil
}

func (st *Storage) EventsForMonth(userID int, date time.Time) ([]model.Event, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	user := st.events[userID]
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	result := []model.Event{}

	for _, eventsByName := range user {
		for _, e := range eventsByName {
			if e.Date.Year() == date.Year() && e.Date.Month() == date.Month() {
				result = append(result, e)
			}
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no events for this month")
	}
	return result, nil
}
