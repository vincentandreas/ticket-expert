package repo

import (
	"context"
	"ticket-expert/models"
)

func (repo *Implementation) SavePromotor(promotor models.Promotor, ctx context.Context) {
	repo.db.WithContext(ctx).Create(&promotor)
}
