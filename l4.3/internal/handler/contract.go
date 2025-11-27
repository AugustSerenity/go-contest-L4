package handler

import (
	"time"

	"github.com/AugustSerenity/go-contest-L4/l4.3_Events-calendar/internal/model"
)

type Service interface {
	CreateEvent(req model.CreateRequest, requestID string) (int, error)
	UpdateEvent(req model.UpdateRequest, requestID string) error
	DeleteEvent(req model.DeleteRequest, requestID string) error

	EventsForDay(userID int, date time.Time) ([]model.Event, error)
	EventsForWeek(userID int, date time.Time) ([]model.Event, error)
	EventsForMonth(userID int, date time.Time) ([]model.Event, error)
}
