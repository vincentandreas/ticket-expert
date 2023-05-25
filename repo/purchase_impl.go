package repo

import "ticket-expert/models"

func (repo *Implementation) SavePurchase(req models.PurchasedTicket) error {
	result := repo.db.Create(&req)
	return result.Error
}
