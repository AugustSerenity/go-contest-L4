package service

import (
	"fmt"
	"time"

	"l2.18/internal/model"
)

type Service struct {
	storage Storage
}

func New(st Storage) *Service {
	return &Service{
		storage: st,
	}
}

func (s *Service) CreateEvent(id int, event model.Event) error {
	if !event.Date.After(time.Now()) {
		return fmt.Errorf("past date")
	}

	if s.storage.ExactEventExists(id, event.Date, event.Name) {
		return fmt.Errorf("event already exists")
	}

	s.storage.Create(id, event)
	return nil
}

func (s *Service) UpdateEvent(userID int, date time.Time, updated model.Event) error {
	if !updated.Date.After(time.Now()) {
		return fmt.Errorf("past date")
	}

	if !updated.Date.Equal(date) {
		return fmt.Errorf("cannot change event date during update")
	}

	if !s.storage.EventAtTimeExists(userID, date) {
		return fmt.Errorf("event does not exist")
	}

	if err := s.storage.Update(userID, date, updated); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteEvent(userID int, date time.Time, name string) error {
	if !s.storage.ExactEventExists(userID, date, name) {
		return fmt.Errorf("event does not exist")
	}

	if err := s.storage.Delete(userID, date, name); err != nil {
		return err
	}

	return nil
}

func (s *Service) ShowEventsForDay(userID int, date time.Time) ([]model.Event, error) {
	res, err := s.storage.EventsForDay(userID, date)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) ShowEventsForWeek(userID int, date time.Time) ([]model.Event, error) {
	res, err := s.storage.EventsForWeek(userID, date)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) ShowEventsForMonth(userID int, date time.Time) ([]model.Event, error) {
	res, err := s.storage.EventsForMonth(userID, date)
	if err != nil {
		return nil, err
	}

	return res, nil

}
