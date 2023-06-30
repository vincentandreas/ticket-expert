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
		bookingData.BookingStatus = "purchased"
		if err := tx.Save(&bookingData).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}
