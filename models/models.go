package models

import (
	"gorm.io/gorm"
)

type Promotor struct {
	gorm.Model
	PromotorName string `json:"promotor_name" validate:"required"`
	Events       []Event
}

type User struct {
	gorm.Model
	FullName         string `json:"full_name" validate:"required" gorm:"not null"`
	UserName         string `json:"user_name" validate:"required" gorm:"not null"`
	Password         string `json:"password" validate:"required" gorm:"not null"`
	BookingTickets   []BookingTicket
	PurchasedTickets []PurchasedTicket
}

type ApiResponse struct {
	Result   string `json:"result"`
	RespCode string `json:"response_code"`
	RespMsg  string `json:"response_message"`
}

type ApiGetResponse struct {
	Result   interface{} `json:"result"`
	RespCode string      `json:"response_code"`
	RespMsg  string      `json:"response_message"`
}

type Event struct {
	gorm.Model
	EventName      string         `json:"event_name" validate:"required"  gorm:"not null"`
	EventCategory  string         `json:"event_category" validate:"required"  gorm:"not null"`
	EventLocation  string         `json:"event_location" validate:"required"  gorm:"not null"`
	PromotorID     uint           `json:"promotor_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	EventDetails   []*EventDetail `json:"event_details" validate:"min=1,dive"`
	BookingTickets []BookingTicket
}

type EventDetail struct {
	gorm.Model
	TicketClass     string `json:"ticket_class" validate:"required"`
	TicketPrice     string `json:"ticket_price" validate:"required"`
	TicketQuota     uint   `json:"ticket_quota" validate:"required"`
	TicketRemaining uint   `json:"ticket_remaining" validate:"required"`
	EventID         uint   `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	BookingDetail   []*BookingDetail
}

type PurchasedTicket struct {
	gorm.Model
	UserID          uint   `json:"user_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	BookingTicketID uint   `json:"booking_ticket_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	PaymentMethod   string `json:"payment_method" validate:"required"`
	PaymentStatus   string `json:"payment_status" validate:"required"`
}

//todo scheduler for expired the booking.
type BookingTicket struct {
	gorm.Model
	TotalPrice      string           `json:"total_price" gorm:"not null"`
	AdminFee        string           `json:"admin_fee" gorm:"not null"`
	UserID          uint             `json:"user_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	BookingDetails  []*BookingDetail `json:"booking_details" validate:"min=1,dive"`
	PurchasedTicket *PurchasedTicket `validate:"omitempty"`
	BookingStatus   string           `json:"booking_status" gorm:"not null"`
	QueueUniqueCode string           `json:"queue_unique_code" validate:"required"`
	EventID         uint             `json:"event_id" validate:"required"`
}

type BookingDetail struct {
	gorm.Model
	EventDetailID   uint   `json:"event_detail_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	Qty             uint   `json:"qty" validate:"required" gorm:"not null"`
	Price           string `json:"price" validate:"required" gorm:"not null"`
	SubTotal        string `gorm:"not null"`
	BookingTicketID uint   `json:"booking_ticket_id" gorm:"not null"`
}

type Qres struct {
	EventID       uint   `json:"event_id"`
	EventName     string `json:"event_name"`
	EventCategory string `json:"event_category"`
	EventLocation string `json:"event_location"`
	PromotorName  string `json:"promotor_name"`
	PromotorID    uint   `json:"promotor_id"`
}
