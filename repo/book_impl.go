package repo

import (
	"errors"
	"gorm.io/gorm"
	"os"
	"strconv"
	"ticket-expert/models"
)

func (repo *Implementation) SaveBooking(req models.BookingTicket) error {
	var grandTotal float64 = 0
	bookDetails := req.BookingDetails

	for i := 0; i < len(bookDetails); i++ {
		price, _ := strconv.ParseFloat(bookDetails[i].Price, 64)

		total := price * float64(bookDetails[i].Qty)
		grandTotal += total
		bookDetails[i].SubTotal = strconv.FormatFloat(total, 'f', -1, 64)
	}
	admEnv := os.Getenv("ADMIN_FEE")
	admFee, _ := strconv.ParseFloat(admEnv, 64)
	grandTotal += admFee
	req.TotalPrice = strconv.FormatFloat(grandTotal, 'f', -1, 64)
	req.AdminFee = admEnv

	err := repo.db.Transaction(func(tx *gorm.DB) error {
		var evData models.Event
		tx.Preload("EventDetails").Where("events.id = ?", req.EventID).Find(&evData)

		evdetails := evData.EventDetails

		for i := 0; i < len(bookDetails); i++ {
			for j := 0; j < len(evdetails); j++ {
				if evdetails[j].ID == bookDetails[i].EventDetailID {
					if evdetails[j].TicketQuota < bookDetails[i].Qty {
						return errors.New("ticket quota not enough")
					} else {
						deductedQuota := evdetails[j].TicketQuota - bookDetails[i].Qty
						err := tx.Model(evdetails).Where("id = ?", evdetails[j].ID).Update("ticket_quota", deductedQuota).Error
						if err != nil {
							return err
						}
					}
					break
				}
			}
		}

		if err := tx.Create(&req).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}
