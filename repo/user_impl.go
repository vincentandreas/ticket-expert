package repo

import (
	"context"
	"ticket-expert/models"
	"ticket-expert/utilities"
)

func (repo *Implementation) SaveUser(user models.User, ctx context.Context) error {
	user.Password = utilities.HashParams(user.Password)
	err := repo.db.WithContext(ctx).Create(&user).Error
	return err
}

func (repo *Implementation) FindUserById(id uint, ctx context.Context) (*models.User, error) {
	var user *models.User
	selSql := "full_name, user_name, phone_number, role"
	err := repo.db.WithContext(ctx).Select(selSql).Where("id = ?", id).First(&user).Error
	return user, err
}

func (repo *Implementation) Login(req models.UserLogin, ctx context.Context) (*models.User, error) {
	var user *models.User
	hashPasswd := utilities.HashParams(req.Password)
	err := repo.db.WithContext(ctx).Where("user_name = ? AND password = ?", req.UserName, hashPasswd).First(&user).Error

	if err != nil || user == nil {
		return nil, err
	}

	return user, nil
}
