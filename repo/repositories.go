package repo

import "ticket-expert/models"

type UserRepository interface {
	SaveUser(user models.User)
}

type PromotorRepository interface {
	SavePromotor(promotor models.Promotor)
}

type EventRepository interface {
	SaveEvent(event models.Event) error
	FindByCondition(location string, category string) ([]models.Qres, error)
	FindByEventId(id string) (models.Event, error)
}

type BookRepository interface {
	SaveBooking(event models.BookingTicket) error
}
type PurchaseRepository interface {
	SavePurchase(event models.PurchasedTicket) error
}

type AllRepository interface {
	UserRepository
	PromotorRepository
	EventRepository
	BookRepository
	PurchaseRepository
}
