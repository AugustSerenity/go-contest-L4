package service

import (
	"time"

	"l2.18/internal/model"
)

type Storage interface {
	Create(id int, event model.Event)
	ExactEventExists(userID int, date time.Time, name string) bool
	Update(userID int, date time.Time, updatedEvent model.Event) error
	EventAtTimeExists(userID int, date time.Time) bool
	Delete(userID int, date time.Time, name string) error
	EventsForDay(userID int, date time.Time) ([]model.Event, error)
	EventsForWeek(userID int, date time.Time) ([]model.Event, error)
	EventsForMonth(userID int, date time.Time) ([]model.Event, error)
}
