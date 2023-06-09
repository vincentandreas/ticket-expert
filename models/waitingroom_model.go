package models

type NewWaitingUser struct {
	UserId      uint   `json:"user_id" validate:"required"`
	EventId     uint   `json:"event_id" validate:"required"`
	QUniqueCode string `json:"q_unique_code"`
}
