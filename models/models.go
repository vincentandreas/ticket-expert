package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName         string `json:"full_name" validate:"required" gorm:"not null"`
	UserName         string `json:"user_name" validate:"required" gorm:"not null,uniqueIndex"`
	Password         string `json:"password" validate:"required" gorm:"not null"`
	PhoneNumber      string `json:"phone_number" validate:"required" gorm:"not null,uniqueIndex"`
	Role             string `json:"role" validate:"required,oneof=USER PROMOTOR" `
	Events           []Event
	BookingTickets   []BookingTicket
	PurchasedTickets []PurchasedTicket
}
type OutUser struct {
	FullName    string `json:"full_name"`
	UserName    string `json:"user_name"`
	PhoneNumber string `json:"phone_number"`
}

func (usr *User) Extract() OutUser {
	ou := OutUser{
		FullName:    usr.FullName,
		UserName:    usr.UserName,
		PhoneNumber: usr.PhoneNumber,
	}
	return ou
}

type UserLogin struct {
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ApiResponse struct {
	Result   string `json:"result"`
	RespCode string `json:"response_code"`
	RespMsg  string `json:"response_message"`
}

type ApiGetResponse struct {
	Data     interface{} `json:"data"`
	RespCode string      `json:"response_code"`
	RespMsg  string      `json:"response_message"`
}

type Event struct {
	gorm.Model
	EventName      string         `json:"event_name" validate:"required"  gorm:"not null"`
	EventDesc      string         `json:"event_desc" validate:"required"`
	EventCategory  string         `json:"event_category" validate:"required"  gorm:"not null"`
	EventLocation  string         `json:"event_location" validate:"required"  gorm:"not null"`
	UserID         uint           `json:"user_id" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	EventDetails   []*EventDetail `json:"event_details" validate:"min=1,dive"`
	BookingTickets []BookingTicket
}

type EventDetail struct {
	gorm.Model
	TicketClass     string `json:"ticket_class" validate:"required"`
	TicketPrice     string `json:"ticket_price" validate:"required"`
	TicketQuota     uint   `json:"ticket_quota" validate:"required"`
	TicketRemaining uint   `json:"ticket_remaining"`
	EventID         uint   `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	BookingDetail   []*BookingDetail
}

type PurchasedTicket struct {
	gorm.Model
	UserID          uint   `validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	BookingTicketID uint   `json:"booking_ticket_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	PaymentMethod   string `json:"payment_method" validate:"required"`
	PaymentStatus   string `json:"payment_status" validate:"required"`
}

//todo scheduler for expired the booking.
type BookingTicket struct {
	gorm.Model
	TotalPrice      string           `json:"total_price" gorm:"not null"`
	AdminFee        string           `json:"admin_fee" gorm:"not null"`
	UserID          uint             `validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	BookingDetails  []*BookingDetail `json:"booking_details" validate:"min=1,dive"`
	PurchasedTicket *PurchasedTicket `validate:"omitempty"`
	BookingStatus   string           `json:"booking_status" gorm:"not null"`
	QUniqueCode     string           `json:"q_unique_code" validate:"required"`
	EventID         uint             `json:"event_id" validate:"required"`
}

type BookingDetail struct {
	gorm.Model
	EventDetailID   uint   `json:"event_detail_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	Qty             uint   `json:"qty" validate:"required" gorm:"not null"`
	Price           string `gorm:"not null"`
	SubTotal        string `gorm:"not null"`
	BookingTicketID uint   `json:"booking_ticket_id" gorm:"not null"`
}

type Qres struct {
	ID            uint   `json:"event_id"`
	EventName     string `json:"event_name"`
	EventCategory string `json:"event_category"`
	EventLocation string `json:"event_location"`
	FullName      string `json:"full_name"`
	UserID        uint   `json:"user_id"`
}

type TicketDetails struct {
}

type PurchaseDetails struct {
	TicketPrice   string `json:"ticket_price"`
	EventCategory string `json:"event_category"`
	EventName     string `json:"event_name"`
	QUniqueCode   string `json:"q_unique_code"`
	BookingStatus string `json:"booking_status"`
	TotalPrice    string `json:"total_price"`
}

type ShowBooking struct {
	EventName     string `json:"event_name"`
	QUniqueCode   string `json:"q_unique_code"`
	BookingStatus string `json:"booking_status"`
	TotalPrice    string `json:"total_price"`
}
