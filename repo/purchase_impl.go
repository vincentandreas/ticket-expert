package repo

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"ticket-expert/models"
)

func (repo *Implementation) SavePurchase(req models.PurchasedTicket, ctx context.Context) error {
	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var bookingData models.BookingTicket
		err := tx.Preload("PurchasedTicket").Where("booking_tickets.id = ?", req.BookingTicketID).Find(&bookingData).Error
		if err != nil {
			return err
		}

		if bookingData.BookingStatus != "active" {
			return errors.New("booking status not active")
		}

		if bookingData.PurchasedTicket != nil {
			return errors.New("this booking already purchased")
		}
		if err := tx.Create(&req).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}

//func (repo *Implementation) FindPurchasedEventById(userId uint, ctx context.Context) {
//	//var
//	selSql := "booking_tickets."
//	repo.db.WithContext(ctx).Preload("BookingTickets").Preload("Events")
//}

func (repo *Implementation) GetTicketDetails(purchasedTicketID uint) (models.TicketDetails, error) {
	var ticketDetails models.TicketDetails
	var purchasedTicket models.PurchasedTicket
	var bookingTicket models.BookingTicket
	var event models.Event
	var eventDetail models.EventDetail

	if err := repo.db.Preload("BookingTickets").First(&purchasedTicket, purchasedTicketID).Error; err != nil {
		return ticketDetails, err
	}

	if err := repo.db.Preload("BookingDetails").First(&bookingTicket, purchasedTicket.BookingTicketID).Error; err != nil {
		return ticketDetails, err
	}

	if err := repo.db.Preload("EventDetails").First(&event, bookingTicket.EventID).Error; err != nil {
		return ticketDetails, err
	}

	if err := repo.db.First(&eventDetail, "event_id = ?", event.ID).Error; err != nil {
		return ticketDetails, err
	}

	ticketDetails.TicketPrice = eventDetail.TicketPrice
	ticketDetails.EventCategory = event.EventCategory
	ticketDetails.EventName = event.EventName
	ticketDetails.QUniqueCode = bookingTicket.QUniqueCode
	ticketDetails.BookingStatus = bookingTicket.BookingStatus
	ticketDetails.TotalPrice = bookingTicket.TotalPrice

	return ticketDetails, nil
}
