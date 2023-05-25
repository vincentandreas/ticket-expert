package repo

import "ticket-expert/models"

func (repo *Implementation) SaveBooking(req models.BookingTicket) error {
	result := repo.db.Create(&req)
	return result.Error
}
