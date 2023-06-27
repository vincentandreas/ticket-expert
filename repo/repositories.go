package repo

import (
	"context"
	"gorm.io/gorm"
	"ticket-expert/models"
)

type UserRepository interface {
	SaveUser(user models.User, ctx context.Context) error
	FindUserById(id uint, ctx context.Context) (*models.User, error)
	Login(req models.UserLogin, ctx context.Context) (*models.User, error)
}

type EventRepository interface {
	SaveEvent(event models.Event, ctx context.Context) error
	FindEventByCondition(eventName string, location string, category string, ctx context.Context) ([]models.Qres, error)
	FindByEventId(id string, ctx context.Context) (models.Event, error)
	FindEvDetailPrice(evDetailIds []uint, ctx context.Context) (map[uint]string, error)
}

type BookRepository interface {
	SaveBooking(event models.BookingTicket, ctx context.Context) error
	UpdTicketQty(id uint, quota uint, tx *gorm.DB, ctx context.Context) error
	CheckBookingPeriodically(ctx context.Context)
	FindBookingByUserId(userId uint, ctx context.Context) ([]*models.ShowBooking, error)
}
type PurchaseRepository interface {
	SavePurchase(event models.PurchasedTicket, ctx context.Context) error
	FindPurchasedEventById(purchasedTicketID string) (models.PurchaseDetails, error)
}

type WaitingRepository interface {
	SaveWaitingQueue(wuser models.NewWaitingUser, ctx context.Context)
	CheckOrderRoom(eventId uint, ctx context.Context) []string
	GetUserInOrderRoom(userId uint, eventId uint, ctx context.Context) string
	CountTotalPeopleInOrderRoom(eventId uint, ctx context.Context) int64
	CountTotalPeopleInWaitingRoom(eventId uint, ctx context.Context) int64
}

type AllRepository interface {
	UserRepository
	EventRepository
	BookRepository
	PurchaseRepository
	WaitingRepository
}
