package model

import "time"

type Request struct {
	UserID    int    `json:"user_id"`
	EventName string `json:"event_name"`
	Date      string `json:"date"`
}

func CastToEvent(requestEvent Request) (Event, error) {
	date, err := time.Parse("2006-01-02 15:04:05", requestEvent.Date)
	if err != nil {
		return Event{}, err
	}
	return Event{Name: requestEvent.EventName, Date: date}, nil
}
