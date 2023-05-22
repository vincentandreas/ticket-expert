package repo

import "ticket-expert/models"

func (repo *Implementation) SaveEvent(event models.Event) {
	repo.db.Create(&event)
}
