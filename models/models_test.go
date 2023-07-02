package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestUser_failed_when_fullname_empty(t *testing.T) {
	req := User{
		UserName:    "qwww",
		Password:    "qwww",
		Role:        "PROMOTOR",
		PhoneNumber: "08120812",
	}

	validationErr := checkStruct(req)
	assert.Contains(t, validationErr[0].Error(), "'FullName' failed")
}

func TestUser_failed_when_role_empty(t *testing.T) {
	req := User{
		UserName:    "qwww",
		Password:    "qwww",
		FullName:    "vvv",
		PhoneNumber: "08120812",
	}

	validationErr := checkStruct(req)
	assert.Contains(t, validationErr[0].Error(), "'Role' failed")
}

func checkStruct(req interface{}) validator.ValidationErrors {
	val := validator.New()
	err := val.Struct(req)
	validationErr := err.(validator.ValidationErrors)
	return validationErr
}

func TestEvent_failed_when_eventdetail_empty(t *testing.T) {
	req := Event{
		EventName:     "qwww",
		EventCategory: "qwww",
		EventLocation: "jakarta",
		EventDesc:     "desc",
		EventPhoto:    "http://www.a.com",
		UserID:        1,
		EventDetails:  nil,
	}

	validationErr := checkStruct(req)
	log.Println("Isi validate Err")
	log.Println(validationErr)
	assert.Contains(t, validationErr[0].Error(), "'EventDetails' failed")
}

func TestEventDetail_failed_when_ticketclass_empty(t *testing.T) {
	req := EventDetail{
		TicketClass:     "",
		TicketPrice:     "234",
		TicketQuota:     2,
		TicketRemaining: 2,
	}

	validationErr := checkStruct(req)
	assert.Contains(t, validationErr[0].Error(), "'TicketClass' failed")
}

//func TestBooking_failed_when_qty_empty(t *testing.T) {
//	req := BookingTicket{
//		UserID:        4,
//		EventDetailID: 3,
//	}
//
//	validationErr := checkStruct(req)
//	assert.Contains(t, validationErr[0].Error(), "'Qty' failed")
//}
