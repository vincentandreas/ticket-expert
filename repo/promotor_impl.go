package repo

import (
	"ticket-expert/models"
)

func (repo *Implementation) SavePromotor(promotor models.Promotor) {
	repo.db.Create(&promotor)
}
