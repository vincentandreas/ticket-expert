package models

import (
	"gorm.io/gorm"
	"time"
)

type Promotor struct {
	gorm.Model
	PromotorName string `json:"promotor_name" validate:"required"`
	Events       []Event
}

//type User struct {
//	gorm.Model
//	Email                string `validate:"omitempty,optional_between=6 75,email" gorm:"type:varchar(75);not null;unique" json:"email"`
//	Password             string `gorm:"not null" json:"-"`
//	Name                 string `validate:"optional_between=1 100" gorm:"type:varchar(100);not null" json:"name"`
//	Role                 int    `validate:"optional_between=1 6" gorm:"not null" json:"role"`
//	Status               int    `validate:"optional_between=1 2" json:"status"`
//	PlainPassword        string `validate:"optional_between=8 50" gorm:"-" json:"plainPassword"`
//	PlainPasswordConfirm string `validate:"optional_between=8 50" gorm:"-" json:"plainPasswordConfirm"`
//}

//
//type GenSiteParam struct {
//	Username   string `json:"username" validate:"required,omitempty"`
//	Password   string `json:"password" validate:"required"`
//	Site       string `json:"site" validate:"required"`
//	KeyCounter int    `json:"keyCounter" validate:"required,min=1,max=4294967295"`
//	KeyPurpose string `json:"keyPurpose" validate:"required,oneof=password loginName answer"`
//	KeyType    string `json:"keyType" validate:"required,oneof=med long max short basic pin name phrase"`
//}
type User struct {
	gorm.Model
	FullName         string `json:"full_name" validate:"required"`
	UserName         string `json:"user_name" validate:"required"`
	Password         string `json:"password" validate:"required"`
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
	EventName     string         `json:"event_name" validate:"required"`
	EventCategory string         `json:"event_category" validate:"required"`
	EventLocation string         `json:"event_location" validate:"required"`
	PromotorID    uint64         `json:"promotor_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	EventDetails  []*EventDetail `json:"event_details" validate:"min=1,dive"`
}

type EventDetail struct {
	gorm.Model
	TicketClass      string `json:"ticket_class" validate:"required"`
	TicketPrice      string `json:"ticket_price" validate:"required"`
	TicketQuota      string `json:"ticket_quota" validate:"required"`
	TicketRemaining  string `json:"ticket_remaining" validate:"required"`
	EventID          uint64 `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	PurchasedTickets []PurchasedTicket
	BookingTickets   []BookingTicket
}

type PurchasedTicket struct {
	gorm.Model
	PurchasedAt     time.Time
	UserID          uint64 `json:"user_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	EventDetailID   uint64 `json:"event_detail_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	BookingTicketID uint64 `json:"booking_ticket_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
}

type BookingTicket struct {
	gorm.Model
	Qty             uint64 `json:"qty" validate:"required"`
	UserID          uint64 `json:"user_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	EventDetailID   uint64 `json:"event_detail_id" validate:"required" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	PurchasedTicket PurchasedTicket
}

type Qres struct {
	EventID       string `json:"event_id"`
	EventName     string `json:"event_name"`
	EventCategory string `json:"event_category"`
	EventLocation string `json:"event_location"`
	PromotorName  string `json:"promotor_name"`
	PromotorID    uint64 `json:"promotor_id"`
}
