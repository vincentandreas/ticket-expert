package models

type NewWaitingUser struct {
	UserId  uint `json:"user_id"`
	EventId uint `json:"event_id"`
}
