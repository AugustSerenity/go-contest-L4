package model

import "time"

type Event struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"event_name"`
	Date      time.Time `json:"date"`
	RemindAt  time.Time `json:"remind_at,omitempty"` 
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"` 
}

type Reminder struct {
	EventID   int       `json:"event_id"`
	UserID    int       `json:"user_id"`
	EventName string    `json:"event_name"`
	RemindAt  time.Time `json:"remind_at"`
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	RequestID string    `json:"request_id,omitempty"`
}
