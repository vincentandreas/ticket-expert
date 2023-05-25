package repo

import "ticket-expert/models"

func (repo *Implementation) SaveEvent(event models.Event) error {
	result := repo.db.Create(&event)
	return result.Error
}
