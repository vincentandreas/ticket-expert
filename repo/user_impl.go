package repo

import (
	"context"
	"gorm.io/gorm"
	"ticket-expert/models"
	"ticket-expert/utilities"
)

func (repo *Implementation) SaveUser(user models.User, ctx context.Context) {
	user.Password = utilities.HashParams(user.Password)
	repo.db.WithContext(ctx).Create(&user)
}

func (repo *Implementation) FindUserById(id uint, ctx context.Context) (*models.User, *gorm.DB) {
	var user *models.User
	res := repo.db.WithContext(ctx).Where("id = ?", id).First(&user)
	return user, res
}
