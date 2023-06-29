package repo

import (
	"context"
	"ticket-expert/models"
)

func (repo *Implementation) SaveEvent(event models.Event, ctx context.Context) error {
	for i := 0; i < len(event.EventDetails); i++ {
		event.EventDetails[i].TicketRemaining = event.EventDetails[i].TicketQuota
	}

	result := repo.db.WithContext(ctx).Create(&event)
	return result.Error
}

func (repo *Implementation) FindEventByCondition(eventName string, location string, category string, ctx context.Context) ([]models.Qres, error) {
	var allres []models.Qres
	selectQ := "events.id, events.event_name, events.event_category, events.event_location, events.event_photo, users.full_name, events.user_id"
	joinQ := "LEFT JOIN users on events.user_id = users.id"
	whereQ := "event_name LIKE ? AND event_location LIKE ? AND event_category LIKE ?"
	result := repo.db.WithContext(ctx).Model(models.Event{}).Select(selectQ).Joins(joinQ).Where(whereQ, "%"+eventName+"%", "%"+location+"%", "%"+category+"%").Scan(&allres)

	return allres, result.Error
}

func (repo *Implementation) FindByEventId(id string, ctx context.Context) (models.Event, error) {
	var events models.Event
	result := repo.db.WithContext(ctx).Preload("EventDetails").Where("id = ?", id).First(&events)
	//result := repo.db.Find()
	return events, result.Error
}

func (repo *Implementation) FindEvDetailPrice(evDetailIds []uint, ctx context.Context) (map[uint]string, error) {
	var events []models.EventDetail
	selCol := "id, ticket_price"
	result := repo.db.WithContext(ctx).Select(selCol).Where("id IN ?", evDetailIds).Find(&events)

	evPrice := make(map[uint]string)

	if events != nil {
		for i := 0; i < len(events); i++ {
			evPrice[events[i].ID] = events[i].TicketPrice
		}
	}
	return evPrice, result.Error
}

func FindEventDetailsByIds(repo *Implementation, ids []uint, ctx context.Context) ([]*models.EventDetail, error) {
	var eventDetails []*models.EventDetail
	result := repo.db.WithContext(ctx).Where("id IN ?", ids).Find(&eventDetails)
	//result := repo.db.Find()
	return eventDetails, result.Error
}
