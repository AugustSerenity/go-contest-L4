package model

type CreateRequest struct {
	UserID   int    `json:"user_id"`
	Name     string `json:"event_name"`
	Date     string `json:"date"`
	RemindIn string `json:"remind_in,omitempty"`
}

type UpdateRequest struct {
	UserID   int    `json:"user_id"`
	Name     string `json:"event_name"`
	Date     string `json:"date"`
	NewName  string `json:"new_name,omitempty"`
	NewDate  string `json:"new_date,omitempty"`
	RemindIn string `json:"remind_in,omitempty"`
}

type DeleteRequest struct {
	UserID int    `json:"user_id"`
	Name   string `json:"event_name"`
	Date   string `json:"date"`
}
