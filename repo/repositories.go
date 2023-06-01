package repo

import (
	"context"
	"gorm.io/gorm"
	"ticket-expert/models"
)

type UserRepository interface {
	SaveUser(user models.User, ctx context.Context)
	FindUserById(id uint, ctx context.Context) (*models.User, error)
}

type PromotorRepository interface {
	SavePromotor(promotor models.Promotor, ctx context.Context)
}

type EventRepository interface {
	SaveEvent(event models.Event, ctx context.Context) error
	FindEventByCondition(location string, category string, ctx context.Context) ([]models.Qres, error)
	FindByEventId(id string, ctx context.Context) (models.Event, error)
}

type BookRepository interface {
	SaveBooking(event models.BookingTicket, ctx context.Context) error
	UpdTicketQty(id uint, quota uint, tx *gorm.DB, ctx context.Context) error
}
type PurchaseRepository interface {
	SavePurchase(event models.PurchasedTicket, ctx context.Context) error
}

type AllRepository interface {
	UserRepository
	PromotorRepository
	EventRepository
	BookRepository
	PurchaseRepository
}
