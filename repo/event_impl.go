package repo

import (
	"context"
	"ticket-expert/models"
)

func (repo *Implementation) SaveEvent(event models.Event, ctx context.Context) error {
	result := repo.db.WithContext(ctx).Create(&event)
	return result.Error
}

func (repo *Implementation) FindEventByCondition(location string, category string, ctx context.Context) ([]models.Qres, error) {

	var allres []models.Qres
	selectQ := "events.id, events.event_name, events.event_category, events.event_location, promotors.promotor_name, events.promotor_id"
	joinQ := "LEFT JOIN promotors on events.promotor_id = promotors.id"
	whereQ := "event_location LIKE ? AND event_category LIKE ?"
	result := repo.db.Model(models.Event{}).Select(selectQ).Joins(joinQ).Where(whereQ, "%"+location+"%", "%"+category+"%").Scan(&allres)

	return allres, result.Error
}

func (repo *Implementation) FindByEventId(id string, ctx context.Context) (models.Event, error) {
	var events models.Event
	result := repo.db.WithContext(ctx).Preload("EventDetails").Where("id = ?", id).Find(&events)
	//result := repo.db.Find()
	return events, result.Error
}

func (repo *Implementation) FindEventDetailsByIds(ids []uint, ctx context.Context) ([]*models.EventDetail, error) {
	var eventDetails []*models.EventDetail
	result := repo.db.WithContext(ctx).Where("id IN ?", ids).Find(&eventDetails)
	//result := repo.db.Find()
	return eventDetails, result.Error
}
