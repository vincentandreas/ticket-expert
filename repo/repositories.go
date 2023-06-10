package repo

import (
	"context"
	"gorm.io/gorm"
	"ticket-expert/models"
)

type UserRepository interface {
	SaveUser(user models.User, ctx context.Context)
	FindUserById(id uint, ctx context.Context) (*models.User, error)
	Login(req models.UserLogin, ctx context.Context) (uint, error)
}

type PromotorRepository interface {
	SavePromotor(promotor models.Promotor, ctx context.Context)
}

type EventRepository interface {
	SaveEvent(event models.Event, ctx context.Context) error
	FindEventByCondition(eventName string, location string, category string, ctx context.Context) ([]models.Qres, error)
	FindByEventId(id string, ctx context.Context) (models.Event, error)
	FindEventDetailsByIds(ids []uint, ctx context.Context) ([]*models.EventDetail, error)
}

type BookRepository interface {
	SaveBooking(event models.BookingTicket, ctx context.Context) error
	UpdTicketQty(id uint, quota uint, tx *gorm.DB, ctx context.Context) error
	CheckBookingPeriod(ctx context.Context)
}
type PurchaseRepository interface {
	SavePurchase(event models.PurchasedTicket, ctx context.Context) error
}

type WaitingRepository interface {
	SaveWaitingQueue(wuser models.NewWaitingUser, ctx context.Context)
	PopWaitingQueue(eventId uint, ctx context.Context) string
	SaveUserInOrderRoom(eventId uint, userIdStr string, qUniqueCode string, ctx context.Context)
	PopUserInOrderRoom(userId uint, eventId uint, ctx context.Context)
	CheckOrderRoom(eventId uint, ctx context.Context) []string
	GetUserInOrderRoom(userId uint, eventId uint, ctx context.Context) string
}

type AllRepository interface {
	UserRepository
	PromotorRepository
	EventRepository
	BookRepository
	PurchaseRepository
	WaitingRepository
}
