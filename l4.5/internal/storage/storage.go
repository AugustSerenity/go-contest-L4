package storage

import (
	"fmt"
	"sync"
	"time"

	"l2.18/internal/model"
)

type Storage struct {
	mu     sync.RWMutex
	events map[int][]model.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[int][]model.Event),
	}
}

func (st *Storage) Create(userID int, event model.Event) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.events[userID] = append(st.events[userID], event)
}

func (st *Storage) Update(userID int, date time.Time, updated model.Event) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	userEvents, ok := st.events[userID]
	if !ok {
		return fmt.Errorf("user not found")
	}

	for i, e := range userEvents {
		if e.Date.Equal(date) {
			st.events[userID][i] = updated
			return nil
		}
	}

	return fmt.Errorf("event not found")
}

func (st *Storage) ExactEventExists(userID int, date time.Time, name string) bool {
	st.mu.RLock()
	defer st.mu.RUnlock()

	userEvents, ok := st.events[userID]
	if !ok {
		return false
	}

	for _, e := range userEvents {
		if e.Date.Equal(date) && e.Name == name {
			return true
		}
	}

	return false
}

func (st *Storage) EventAtTimeExists(userID int, date time.Time) bool {
	st.mu.RLock()
	defer st.mu.RUnlock()

	userEvents, ok := st.events[userID]
	if !ok {
		return false
	}

	for _, e := range userEvents {
		if e.Date.Equal(date) {
			return true
		}
	}

	return false
}

func (st *Storage) Delete(userID int, date time.Time, name string) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	userEvents, ok := st.events[userID]
	if !ok {
		return fmt.Errorf("user not found")
	}

	for i, e := range userEvents {
		if e.Date.Equal(date) && e.Name == name {
			st.events[userID] = append(userEvents[:i], userEvents[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("event not found")
}

func (st *Storage) EventsForDay(userID int, date time.Time) ([]model.Event, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	userEvents, ok := st.events[userID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	var result []model.Event
	for _, event := range userEvents {
		if event.Date.Year() == date.Year() &&
			event.Date.Month() == date.Month() &&
			event.Date.Day() == date.Day() {
			result = append(result, event)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no events for this day")
	}

	return result, nil
}

func (st *Storage) EventsForWeek(userID int, date time.Time) ([]model.Event, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	userEvents, ok := st.events[userID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	var result []model.Event
	targetYear, targetWeek := date.ISOWeek()
	for _, event := range userEvents {
		eventYear, eventWeek := event.Date.ISOWeek()
		if eventYear == targetYear &&
			eventWeek == targetWeek {
			result = append(result, event)
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

	userEvents, ok := st.events[userID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	var result []model.Event
	for _, event := range userEvents {
		if event.Date.Year() == date.Year() &&
			event.Date.Month() == date.Month() {
			result = append(result, event)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no events for this month")
	}

	return result, nil
}
