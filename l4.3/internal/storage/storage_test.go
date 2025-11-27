package storage

import (
	"testing"
	"time"

	"github.com/AugustSerenity/go-contest-L4/l4.3_Events-calendar/internal/model"
)

func newEvent(userID int, name string, date time.Time) model.Event {
	return model.Event{
		UserID: userID,
		Name:   name,
		Date:   date,
	}
}

func TestCreate(t *testing.T) {
	st := New()
	date := time.Now().Add(24 * time.Hour)

	id := st.Create(newEvent(1, "test", date))

	if id != 1 {
		t.Fatalf("expected id=1, got %d", id)
	}

	if len(st.events[1]) != 1 {
		t.Fatalf("expected 1 event in storage")
	}

	if !st.events[1][0].IsActive {
		t.Fatalf("new event must be active")
	}
}

func TestUpdate(t *testing.T) {
	st := New()
	date := time.Now().Add(24 * time.Hour)

	st.Create(newEvent(1, "old", date))

	updated := newEvent(1, "newName", date)
	err := st.Update(1, date, updated)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ev := st.events[1][0]
	if ev.Name != "newName" {
		t.Fatalf("expected updated name newName, got %s", ev.Name)
	}
}

func TestUpdate_EventNotFound(t *testing.T) {
	st := New()
	date := time.Now().Add(24 * time.Hour)

	err := st.Update(1, date, newEvent(1, "abc", date))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetEventByTime(t *testing.T) {
	st := New()
	date := time.Now().Add(24 * time.Hour)

	st.Create(newEvent(1, "meeting", date))

	ev, err := st.GetEventByTime(1, date, "meeting")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ev.Name != "meeting" {
		t.Fatalf("expected meeting, got %s", ev.Name)
	}
}

func TestGetEventByTime_NotFound(t *testing.T) {
	st := New()

	_, err := st.GetEventByTime(1, time.Now(), "none")
	if err == nil {
		t.Fatalf("expected error for missing event")
	}
}

func TestExactEventExists(t *testing.T) {
	st := New()
	date := time.Now().Add(24 * time.Hour)

	st.Create(newEvent(1, "ping", date))

	if !st.ExactEventExists(1, date, "ping") {
		t.Fatalf("event should exist")
	}

	if st.ExactEventExists(1, date, "pong") {
		t.Fatalf("event should NOT exist")
	}
}

func TestDelete(t *testing.T) {
	st := New()
	date := time.Now().Add(24 * time.Hour)

	st.Create(newEvent(1, "meet", date))

	err := st.Delete(1, date, "meet")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if st.events[1][0].IsActive {
		t.Fatalf("expected event to be inactive after delete")
	}
}

func TestDelete_NotFound(t *testing.T) {
	st := New()
	date := time.Now()

	err := st.Delete(1, date, "absent")
	if err == nil {
		t.Fatalf("expected error for missing event")
	}
}

func TestGetEventsForDay(t *testing.T) {
	st := New()
	date := time.Date(2025, 12, 5, 10, 0, 0, 0, time.UTC)

	st.Create(newEvent(1, "one", date))
	st.Create(newEvent(1, "two", date.Add(2*time.Hour)))

	events, err := st.GetEventsForDay(1, date)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestGetEventsForDay_NoEvents(t *testing.T) {
	st := New()

	_, err := st.GetEventsForDay(1, time.Now())
	if err == nil {
		t.Fatalf("expected error for no events")
	}
}

func TestGetEventsForWeek(t *testing.T) {
	st := New()

	// same ISO week
	date1 := time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC)
	date2 := date1.Add(48 * time.Hour)

	st.Create(newEvent(1, "a", date1))
	st.Create(newEvent(1, "b", date2))

	events, err := st.GetEventsForWeek(1, date1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestGetEventsForWeek_NoEvents(t *testing.T) {
	st := New()

	_, err := st.GetEventsForWeek(1, time.Now())
	if err == nil {
		t.Fatalf("expected error for empty week")
	}
}

func TestGetEventsForMonth(t *testing.T) {
	st := New()

	date1 := time.Date(2025, 10, 4, 10, 0, 0, 0, time.UTC)
	date2 := time.Date(2025, 10, 15, 10, 0, 0, 0, time.UTC)

	st.Create(newEvent(1, "a", date1))
	st.Create(newEvent(1, "b", date2))

	events, err := st.GetEventsForMonth(1, date1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestGetEventsForMonth_NoEvents(t *testing.T) {
	st := New()

	_, err := st.GetEventsForMonth(1, time.Now())
	if err == nil {
		t.Fatalf("expected error for empty month")
	}
}

func TestArchiveOldEvents(t *testing.T) {
	st := New()

	oldDate := time.Now().Add(-48 * time.Hour)
	newDate := time.Now().Add(48 * time.Hour)

	st.Create(newEvent(1, "old1", oldDate))
	st.Create(newEvent(1, "old2", oldDate))
	st.Create(newEvent(1, "new", newDate))

	cutoff := time.Now().Add(-24 * time.Hour)
	n := st.ArchiveOldEvents(cutoff)

	if n != 2 {
		t.Fatalf("expected to archive 2 events, got %d", n)
	}

	if st.events[1][0].IsActive || st.events[1][1].IsActive {
		t.Fatalf("old events should be inactive")
	}

	if !st.events[1][2].IsActive {
		t.Fatalf("new event should remain active")
	}
}
