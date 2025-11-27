package worker

import (
	"fmt"
	"time"

	"github.com/AugustSerenity/go-contest-L4/l4.3_Events-calendar/internal/model"
)

type ReminderWorker struct {
	reminderChan chan model.Reminder
	done         chan bool
}

func New() *ReminderWorker {
	return &ReminderWorker{
		reminderChan: make(chan model.Reminder, 100),
		done:         make(chan bool),
	}
}

func (w *ReminderWorker) Start() {
	go func() {
		for reminder := range w.reminderChan {
			duration := time.Until(reminder.RemindAt)
			if duration > 0 {
				time.Sleep(duration)
				w.sendReminder(reminder)
			}
		}
		w.done <- true
	}()
}

func (w *ReminderWorker) Stop() {
	close(w.reminderChan)
	<-w.done
}

func (w *ReminderWorker) AddReminder(reminder model.Reminder) {
	select {
	case w.reminderChan <- reminder:
	default:
		fmt.Printf("Reminder channel full, dropping reminder for event %d\n", reminder.EventID)
	}
}

func (w *ReminderWorker) sendReminder(reminder model.Reminder) {
	fmt.Printf("REMINDER: User %d - Event '%s' is scheduled!\n",
		reminder.UserID, reminder.EventName)
}
