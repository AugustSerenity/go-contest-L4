package service

import (
	"testing"
	"time"

	"l2.18/internal/model"
	"l2.18/internal/storage"
)

func TestCreateEvent_PastDate(t *testing.T) {
	st := storage.New()
	srv := New(st)

	event := model.Event{
		Name: "past event",
		Date: time.Now().Add(-1 * time.Hour),
	}

	err := srv.CreateEvent(1, event)
	if err == nil {
		t.Fatalf("expected error for past date, got nil")
	}
}

func TestCreateEvent_Duplicate(t *testing.T) {
	st := storage.New()
	srv := New(st)

	event := model.Event{
		Name: "duplicate",
		Date: time.Now().Add(24 * time.Hour),
	}

	err := srv.CreateEvent(1, event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = srv.CreateEvent(1, event)
	if err == nil {
		t.Fatalf("expected error for duplicate event, got nil")
	}
}

func TestCreateEvent_Success(t *testing.T) {
	st := storage.New()
	srv := New(st)

	event := model.Event{
		Name: "new event",
		Date: time.Now().Add(2 * time.Hour),
	}

	err := srv.CreateEvent(1, event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateEvent_ChangeDateNotAllowed(t *testing.T) {
	st := storage.New()
	srv := New(st)

	originalDate := time.Now().Add(24 * time.Hour)
	event := model.Event{Name: "event1", Date: originalDate}
	st.Create(1, event)

	updated := model.Event{Name: "updated", Date: originalDate.Add(1 * time.Hour)}

	err := srv.UpdateEvent(1, originalDate, updated)
	if err == nil {
		t.Fatalf("expected error for changing event date, got nil")
	}
}

func TestUpdateEvent_NotExist(t *testing.T) {
	st := storage.New()
	srv := New(st)

	originalDate := time.Now().Add(24 * time.Hour)
	updated := model.Event{Name: "updated", Date: originalDate}

	err := srv.UpdateEvent(1, originalDate, updated)
	if err == nil {
		t.Fatalf("expected error for non-existent event, got nil")
	}
}

func TestUpdateEvent_Success(t *testing.T) {
	st := storage.New()
	srv := New(st)

	originalDate := time.Now().Add(24 * time.Hour)
	event := model.Event{Name: "event1", Date: originalDate}
	st.Create(1, event)

	updated := model.Event{Name: "updated", Date: originalDate}
	err := srv.UpdateEvent(1, originalDate, updated)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteEvent_NotExist(t *testing.T) {
	st := storage.New()
	srv := New(st)

	date := time.Now().Add(24 * time.Hour)
	err := srv.DeleteEvent(1, date, "nonexistent")
	if err == nil {
		t.Fatalf("expected error deleting non-existent event, got nil")
	}
}

func TestDeleteEvent_Success(t *testing.T) {
	st := storage.New()
	srv := New(st)

	date := time.Now().Add(24 * time.Hour)
	event := model.Event{Name: "to delete", Date: date}
	st.Create(1, event)

	err := srv.DeleteEvent(1, date, "to delete")
	if err != nil {
		t.Fatalf("unexpected error deleting event: %v", err)
	}
}

func TestUpdateEvent_PastDate(t *testing.T) {
	st := storage.New()
	srv := New(st)

	originalDate := time.Now().Add(24 * time.Hour)
	event := model.Event{Name: "event1", Date: originalDate}
	st.Create(1, event)

	updated := model.Event{Name: "updated", Date: time.Now().Add(-1 * time.Hour)}

	err := srv.UpdateEvent(1, originalDate, updated)
	if err == nil {
		t.Fatalf("expected error for updating with past date, got nil")
	}
}

func TestDeleteEvent_UserNotExist(t *testing.T) {
	st := storage.New()
	srv := New(st)

	date := time.Now().Add(24 * time.Hour)
	err := srv.DeleteEvent(999, date, "event")
	if err == nil {
		t.Fatalf("expected error deleting event for non-existent user, got nil")
	}
}

func TestCreateEvent_EmptyName(t *testing.T) {
	st := storage.New()
	srv := New(st)

	event := model.Event{
		Name: "",
		Date: time.Now().Add(2 * time.Hour),
	}

	err := srv.CreateEvent(1, event)
	if err != nil {
		t.Fatalf("unexpected error creating event with empty name: %v", err)
	}
}

func TestUpdateEvent_EmptyName(t *testing.T) {
	st := storage.New()
	srv := New(st)

	date := time.Now().Add(24 * time.Hour)
	event := model.Event{Name: "test", Date: date}
	st.Create(1, event)

	updated := model.Event{Name: "", Date: date}
	err := srv.UpdateEvent(1, date, updated)
	if err != nil {
		t.Fatalf("unexpected error updating event with empty name: %v", err)
	}
}
