package repo

import (
	"ticket-expert/models"
	"ticket-expert/utilities"
)

func (repo *Implementation) SaveUser(user models.User) {
	user.Password = utilities.HashParams(user.Password)
	repo.db.Create(&user)
}
