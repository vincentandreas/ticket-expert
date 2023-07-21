package repo

import (
	"ticket-expert/models"
)

/// siapa punya duit bagi dong, dari masa lalu ini. .
func (repo *Implementation) SavePromotor(promotor models.Promotor) {
	repo.db.Create(&promotor)
}
