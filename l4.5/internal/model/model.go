package model

import "time"

type Event struct {
	Name string    `json:"event_name"`
	Date time.Time `json:"date"`
}
