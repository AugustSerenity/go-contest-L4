package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/AugustSerenity/go-contest-L4/l4.3_Events-calendar/internal/model"
)

type MemoryStorage struct {
	mu           sync.RWMutex
	events       map[int][]model.Event
	eventCounter int
}

func New() *MemoryStorage {
	return &MemoryStorage{
		events:       make(map[int][]model.Event),
		eventCounter: 1,
	}
}

func (st *MemoryStorage) Create(event model.Event) int {
	st.mu.Lock()
	defer st.mu.Unlock()

	event.ID = st.eventCounter
	event.CreatedAt = time.Now()
	event.IsActive = true
	st.eventCounter++

	st.events[event.UserID] = append(st.events[event.UserID], event)
	return event.ID
}

func (st *MemoryStorage) Update(userID int, date time.Time, updated model.Event) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	userEvents, ok := st.events[userID]
	if !ok {
		return fmt.Errorf("user not found")
	}

	for i, e := range userEvents {
		if e.Date.Equal(date) && e.IsActive {
			updated.ID = e.ID
			updated.CreatedAt = e.CreatedAt
			updated.IsActive = true
			st.events[userID][i] = updated
			return nil
		}
	}

	return fmt.Errorf("event not found")
}

func (st *MemoryStorage) GetEventByTime(userID int, date time.Time, name string) (*model.Event, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	userEvents, ok := st.events[userID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	for i := range userEvents {
		e := &userEvents[i]
		if e.Date.Equal(date) && e.Name == name && e.IsActive {
			return e, nil
		}
	}

	return nil, fmt.Errorf("event not found")
}

func (st *MemoryStorage) ExactEventExists(userID int, date time.Time, name string) bool {
	st.mu.RLock()
	defer st.mu.RUnlock()

	events, ok := st.events[userID]
	if !ok {
		return false
	}

	for _, e := range events {
		if e.Name == name && e.Date.Equal(date) && e.IsActive {
			return true
		}
	}
	return false
}

func (st *MemoryStorage) Delete(userID int, date time.Time, name string) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	events, ok := st.events[userID]
	if !ok {
		return fmt.Errorf("user not found")
	}

	for i := range events {
		if events[i].Name == name && events[i].Date.Equal(date) && events[i].IsActive {
			st.events[userID][i].IsActive = false
			return nil
		}
	}

	return fmt.Errorf("event not found")
}

func (st *MemoryStorage) GetEventsForDay(userID int, date time.Time) ([]model.Event, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	events, ok := st.events[userID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	var res []model.Event
	for _, e := range events {
		if e.IsActive &&
			e.Date.Year() == date.Year() &&
			e.Date.Month() == date.Month() &&
			e.Date.Day() == date.Day() {
			res = append(res, e)
		}
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("no events")
	}
	return res, nil
}

func (st *MemoryStorage) GetEventsForWeek(userID int, date time.Time) ([]model.Event, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	events, ok := st.events[userID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	targetY, targetW := date.ISOWeek()

	var res []model.Event
	for _, e := range events {
		y, w := e.Date.ISOWeek()
		if e.IsActive && y == targetY && w == targetW {
			res = append(res, e)
		}
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("no events")
	}
	return res, nil
}

func (st *MemoryStorage) GetEventsForMonth(userID int, date time.Time) ([]model.Event, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	events, ok := st.events[userID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	var res []model.Event
	for _, e := range events {
		if e.IsActive &&
			e.Date.Year() == date.Year() &&
			e.Date.Month() == date.Month() {
			res = append(res, e)
		}
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("no events")
	}
	return res, nil
}

func (st *MemoryStorage) ArchiveOldEvents(cutoff time.Time) int {
	st.mu.Lock()
	defer st.mu.Unlock()

	count := 0

	for userID, arr := range st.events {
		for i := range arr {
			if arr[i].IsActive && arr[i].Date.Before(cutoff) {
				st.events[userID][i].IsActive = false
				count++
			}
		}
	}

	return count
}
