package repo

import "ticket-expert/models"

type UserRepository interface {
	SaveUser(user models.User)
}

type PromotorRepository interface {
	SavePromotor(promotor models.Promotor)
}

type EventRepository interface {
	SaveEvent(event models.Event)
}

type AllRepository interface {
	UserRepository
	PromotorRepository
	EventRepository
}
