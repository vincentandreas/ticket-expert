package repo

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"ticket-expert/models"
)

func (repo *Implementation) SavePurchase(req models.PurchasedTicket) error {
	fmt.Println("------------------------------------------")
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		var bookingData models.BookingTicket
		err := repo.db.Preload("PurchasedTicket").Where("booking_tickets.id = ?", req.BookingTicketID).Find(&bookingData).Error
		if err != nil {
			return err
		}

		if bookingData.BookingStatus != "active" {
			return errors.New("booking status not active")
		}

		if bookingData.PurchasedTicket != nil {
			return errors.New("this booking already purchased")
		}
		if err := repo.db.Create(&req).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}
