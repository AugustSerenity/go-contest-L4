package service

import (
	"time"

	"github.com/AugustSerenity/go-contest-L4/l4.3_Events-calendar/internal/model"
)

type Storage interface {
	Create(event model.Event) int
	Update(userID int, date time.Time, updated model.Event) error
	GetEventByTime(userID int, date time.Time, name string) (*model.Event, error)
	Delete(userID int, date time.Time, name string) error
	GetEventsForDay(userID int, date time.Time) ([]model.Event, error)
	GetEventsForWeek(userID int, date time.Time) ([]model.Event, error)
	GetEventsForMonth(userID int, date time.Time) ([]model.Event, error)
	ArchiveOldEvents(cutoff time.Time) int
	ExactEventExists(userID int, date time.Time, name string) bool
}
