package service

import (
	"fmt"
	"time"

	"github.com/AugustSerenity/go-contest-L4/l4.3_Events-calendar/internal/model"
	"github.com/AugustSerenity/go-contest-L4/l4.3_Events-calendar/internal/service/logger"
	"github.com/AugustSerenity/go-contest-L4/l4.3_Events-calendar/internal/service/worker"
	"github.com/AugustSerenity/go-contest-L4/l4.3_Events-calendar/internal/storage"
)

type Service struct {
	storage        storage.Storage
	logger         *logger.Logger
	reminderWorker *worker.ReminderWorker
	archiveWorker  *worker.ArchiveWorker
}

func New(st storage.Storage) *Service {
	svc := &Service{
		storage:        st,
		logger:         logger.New(),
		reminderWorker: worker.New(),
		archiveWorker:  worker.NewArchiveWorker(st),
	}
	svc.logger.Start()
	svc.reminderWorker.Start()
	svc.archiveWorker.Start()
	return svc
}

func (s *Service) Stop() {
	s.archiveWorker.Stop()
	s.reminderWorker.Stop()
	s.logger.Stop()
}

func parseDay(dateStr string) (time.Time, error) {
	d, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local), nil
}

func (s *Service) CreateEvent(req model.CreateRequest, requestID string) (int, error) {
	s.logger.Info(fmt.Sprintf("Creating event for user %d", req.UserID), requestID)

	eventDate, err := parseDay(req.Date)
	if err != nil {
		return 0, fmt.Errorf("invalid date format (expected YYYY-MM-DD)")
	}

	if !eventDate.After(time.Now().AddDate(0, 0, -1)) {
		return 0, fmt.Errorf("past date")
	}

	if s.storage.ExactEventExists(req.UserID, eventDate, req.Name) {
		return 0, fmt.Errorf("event already exists")
	}

	event := model.Event{
		UserID: req.UserID,
		Name:   req.Name,
		Date:   eventDate,
	}

	if req.RemindIn != "" {
		duration, err := time.ParseDuration(req.RemindIn)
		if err != nil {
			return 0, fmt.Errorf("invalid remind_in format")
		}
		event.RemindAt = eventDate.Add(-duration)
	}

	eventID := s.storage.Create(event)

	if !event.RemindAt.IsZero() {
		s.reminderWorker.AddReminder(model.Reminder{
			EventID:   eventID,
			UserID:    req.UserID,
			EventName: req.Name,
			RemindAt:  event.RemindAt,
		})
	}

	s.logger.Info(fmt.Sprintf("Event created with ID %d", eventID), requestID)
	return eventID, nil
}

func (s *Service) UpdateEvent(req model.UpdateRequest, requestID string) error {
	s.logger.Info("Updating event", requestID)

	date, err := parseDay(req.Date)
	if err != nil {
		return fmt.Errorf("invalid date")
	}

	event, err := s.storage.GetEventByTime(req.UserID, date, req.Name)
	if err != nil {
		return err
	}

	newEvent := *event

	if req.NewName != "" {
		newEvent.Name = req.NewName
	}
	if req.NewDate != "" {
		d, err := parseDay(req.NewDate)
		if err != nil {
			return fmt.Errorf("invalid new_date")
		}
		newEvent.Date = d
	}
	if req.RemindIn != "" {
		dur, err := time.ParseDuration(req.RemindIn)
		if err != nil {
			return fmt.Errorf("invalid remind_in")
		}
		newEvent.RemindAt = newEvent.Date.Add(-dur)
	}

	err = s.storage.Update(req.UserID, date, newEvent)
	if err != nil {
		return err
	}

	s.logger.Info("Event updated", requestID)

	if !newEvent.RemindAt.IsZero() {
		s.reminderWorker.AddReminder(model.Reminder{
			EventID:   newEvent.ID,
			UserID:    newEvent.UserID,
			EventName: newEvent.Name,
			RemindAt:  newEvent.RemindAt,
		})
	}

	return nil
}

func (s *Service) DeleteEvent(req model.DeleteRequest, requestID string) error {
	s.logger.Info("Deleting event", requestID)

	date, err := parseDay(req.Date)
	if err != nil {
		return fmt.Errorf("invalid date")
	}

	err = s.storage.Delete(req.UserID, date, req.Name)
	if err != nil {
		return err
	}

	s.logger.Info("Event deleted", requestID)
	return nil
}

func (s *Service) EventsForDay(userID int, date time.Time) ([]model.Event, error) {
	date = date.Truncate(24 * time.Hour)
	return s.storage.GetEventsForDay(userID, date)
}

func (s *Service) EventsForWeek(userID int, date time.Time) ([]model.Event, error) {
	date = date.Truncate(24 * time.Hour)
	return s.storage.GetEventsForWeek(userID, date)
}

func (s *Service) EventsForMonth(userID int, date time.Time) ([]model.Event, error) {
	date = date.Truncate(24 * time.Hour)
	return s.storage.GetEventsForMonth(userID, date)
}
