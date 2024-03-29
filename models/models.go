package models

import (
	"time"

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

type CustomModel struct {
	ID        uint `gorm:"primarykey" json:"-"` 
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time`json:"-"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
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
	UserID         uint           `json:"creator_id" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	EventPhoto     string         `json:"event_photo" validate:"required"`
	EventDetails   []*EventDetail `json:"event_details" validate:"min=1,dive"`
	BookingTickets []BookingTicket
}

type OutEvent struct {
	ID            uint           `json:"id"`
	EventName     string         `json:"event_name"`
	EventDesc     string         `json:"event_desc"`
	EventCategory string         `json:"event_category"`
	EventLocation string         `json:"event_location"`
	UserID        uint           `json:"creator_id"`
	EventPhoto    string         `json:"event_photo"`
	EventDetails  []*EventDetail `json:"event_details"`
}

func (ev *Event) Extract() OutEvent {
	return OutEvent{
		EventName:     ev.EventName,
		EventDesc:     ev.EventDesc,
		EventCategory: ev.EventCategory,
		EventLocation: ev.EventLocation,
		UserID:        ev.UserID,
		EventPhoto:    ev.EventPhoto,
		EventDetails:  ev.EventDetails,
	}
}

type EventDetail struct {
	CustomModel		
	TicketClass     string `json:"ticket_class" validate:"required"`
	TicketPrice     string `json:"ticket_price" validate:"required"`
	TicketQuota     uint   `json:"ticket_quota" validate:"required"`
	TicketRemaining uint   `json:"ticket_remaining"`
	EventID         uint   `json:"-" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	BookingDetail   []*BookingDetail `json:"-"`
}

type PurchasedTicket struct {
	gorm.Model		
	UserID          uint   `validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	BookingTicketID uint   `json:"booking_ticket_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	PaymentMethod   string `json:"payment_method" validate:"required"`
	TransRefNo      string `json:"trans_ref_no"  validate:"required"`
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
	EventPhoto    string `json:"event_photo"`
	FullName      string `json:"full_name"`
	UserID        uint   `json:"user_id"`
}

type TicketDetails struct {
	Qty         uint   `json:"qty"`
	Price       string `json:"price"`
	SubTotal    string `json:"sub_total"`
	TicketClass string `json:"ticket_class"`
}

type PurchaseDetails struct {
	AdminFee      string          `json:"admin_fee"`
	EventCategory string          `json:"event_category"`
	EventName     string          `json:"event_name"`
	QUniqueCode   string          `json:"q_unique_code"`
	BookingStatus string          `json:"booking_status"`
	TotalPrice    string          `json:"total_price"`
	TicketDetails []TicketDetails `json:"ticket_details"`
}

type ShowBooking struct {
	EventName     string `json:"event_name"`
	QUniqueCode   string `json:"q_unique_code"`
	BookingStatus string `json:"booking_status"`
	TotalPrice    string `json:"total_price"`
}

type PurchaseReq struct {
	BookingUniqCode string `json:"booking_uniq_code" validate:"required"`
	PaymentMethod   string `json:"payment_method" validate:"required"`
	TransRefNo      string `json:"trans_ref_no"  validate:"required"`
}
