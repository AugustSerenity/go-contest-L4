package handler

import (
	"time"

	"l2.18/internal/model"
)

type Service interface {
	CreateEvent(id int, event model.Event) error
	UpdateEvent(userID int, date time.Time, updated model.Event) error
	DeleteEvent(userID int, date time.Time, name string) error
	ShowEventsForDay(userID int, date time.Time) ([]model.Event, error)
	ShowEventsForWeek(userID int, date time.Time) ([]model.Event, error)
	ShowEventsForMonth(userID int, date time.Time) ([]model.Event, error)
}
