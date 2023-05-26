package repo

import "ticket-expert/models"

func (repo *Implementation) SaveEvent(event models.Event) error {
	result := repo.db.Create(&event)
	return result.Error
}

func (repo *Implementation) FindByCondition(location string, category string) ([]models.Qres, error) {

	var allres []models.Qres
	selectQ := "events.id, events.event_name, events.event_category, events.event_location, promotors.promotor_name, events.promotor_id"
	joinQ := "left join promotors on events.promotor_id = promotors.id"
	whereQ := "event_location LIKE ? AND event_category LIKE ?"
	result := repo.db.Model(models.Event{}).Select(selectQ).Joins(joinQ).Where(whereQ, "%"+location+"%", "%"+category+"%").Scan(&allres)

	return allres, result.Error
}

func (repo *Implementation) FindByEventId(id string) (models.Event, error) {
	var events models.Event
	result := repo.db.Preload("EventDetails").Where("id = ?", id).Find(&events)
	//result := repo.db.Find()
	return events, result.Error
}
