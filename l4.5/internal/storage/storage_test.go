package storage

import (
	"testing"
	"time"

	"l2.18/internal/model"
)

func TestCreate(t *testing.T) {
	st := New()
	event := model.Event{
		Name: "event1",
		Date: time.Now().Add(24 * time.Hour),
	}

	st.Create(1, event)
	events, err := st.EventsForDay(1, event.Date)
	if err != nil {
		t.Fatalf("expected event, got error: %v", err)
	}
	if len(events) != 1 || events[0].Name != event.Name {
		t.Fatalf("expected event to be stored correctly")
	}
}

func TestUpdate(t *testing.T) {
	st := New()
	userID := 1
	oldDate := time.Now().Add(24 * time.Hour)
	event := model.Event{Name: "event1", Date: oldDate}
	st.Create(userID, event)

	updatedEvent := model.Event{Name: "updated_event", Date: oldDate}
	err := st.Update(userID, oldDate, updatedEvent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events, err := st.EventsForDay(userID, oldDate)
	if err != nil {
		t.Fatalf("expected event after update, got error: %v", err)
	}
	if len(events) != 1 || events[0].Name != "updated_event" {
		t.Fatalf("event not updated properly")
	}
}

func TestDelete(t *testing.T) {
	st := New()
	userID := 1
	date := time.Now().Add(24 * time.Hour)
	event := model.Event{Name: "event1", Date: date}
	st.Create(userID, event)

	err := st.Delete(userID, date, "event1")
	if err != nil {
		t.Fatalf("unexpected error deleting event: %v", err)
	}

	_, err = st.EventsForDay(userID, date)
	if err == nil {
		t.Fatalf("expected error after deletion, got nil")
	}
}

func TestExactEventExists(t *testing.T) {
	st := New()
	userID := 1
	date := time.Now().Add(24 * time.Hour)
	event := model.Event{Name: "event1", Date: date}
	st.Create(userID, event)

	if !st.ExactEventExists(userID, date, "event1") {
		t.Fatalf("expected event to exist")
	}

	if st.ExactEventExists(userID, date, "nonexistent") {
		t.Fatalf("did not expect event to exist")
	}
}

func TestEventAtTimeExists(t *testing.T) {
	st := New()
	userID := 1
	date := time.Now().Add(24 * time.Hour)
	event := model.Event{Name: "event1", Date: date}
	st.Create(userID, event)

	if !st.EventAtTimeExists(userID, date) {
		t.Fatalf("expected event at time to exist")
	}

	if st.EventAtTimeExists(userID, time.Now()) {
		t.Fatalf("did not expect event at this time to exist")
	}
}

func TestEventsForDay(t *testing.T) {
	st := New()
	userID := 1
	date := time.Now().Add(24 * time.Hour)
	event := model.Event{Name: "event1", Date: date}
	st.Create(userID, event)

	events, err := st.EventsForDay(userID, date)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) == 0 {
		t.Fatalf("expected events for day")
	}
}

func TestEventsForWeek(t *testing.T) {
	st := New()
	userID := 1
	date := time.Now().AddDate(0, 0, 1)
	event := model.Event{Name: "event1", Date: date}
	st.Create(userID, event)

	events, err := st.EventsForWeek(userID, date)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) == 0 {
		t.Fatalf("expected events for week")
	}
}

func TestEventsForMonth(t *testing.T) {
	st := New()
	userID := 1
	date := time.Now().AddDate(0, 0, 1)
	event := model.Event{Name: "event1", Date: date}
	st.Create(userID, event)

	events, err := st.EventsForMonth(userID, date)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) == 0 {
		t.Fatalf("expected events for month")
	}
}
