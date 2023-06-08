package repo

import (
	"context"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"ticket-expert/models"
)

func (repo *Implementation) UpdTicketQty(id uint, quota uint, tx *gorm.DB, ctx context.Context) error {
	var evdetails models.EventDetail
	err := tx.WithContext(ctx).Model(evdetails).Where("id = ?", id).Update("ticket_quota", quota).Error
	return err
}

func (repo *Implementation) SaveBooking(req models.BookingTicket, ctx context.Context) error {
	var grandTotal float64 = 0
	bookDetails := req.BookingDetails

	if !repo.isValidUniqueId(req, ctx) {
		return errors.New("Failed, Queue Unique Id Not Match")
	}

	for i := 0; i < len(bookDetails); i++ {
		price, _ := strconv.ParseFloat(bookDetails[i].Price, 64)

		total := price * float64(bookDetails[i].Qty)
		grandTotal += total
		bookDetails[i].SubTotal = strconv.FormatFloat(total, 'f', -1, 64)
	}
	admEnv := os.Getenv("admin_fee")
	admFee, _ := strconv.ParseFloat(admEnv, 64)
	grandTotal += admFee
	req.TotalPrice = strconv.FormatFloat(grandTotal, 'f', -1, 64)
	req.AdminFee = admEnv

	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var evData models.Event
		tx.Preload("EventDetails").Where("events.id = ?", req.EventID).Find(&evData)

		evdetails := evData.EventDetails

		for i := 0; i < len(bookDetails); i++ {
			for j := 0; j < len(evdetails); j++ {
				if evdetails[j].ID == bookDetails[i].EventDetailID {
					if evdetails[j].TicketQuota < bookDetails[i].Qty {
						return errors.New("ticket quota not enough")
					} else {
						deductedQuota := evdetails[j].TicketQuota - bookDetails[i].Qty
						err := repo.UpdTicketQty(evdetails[j].ID, deductedQuota, tx, ctx)
						//err := tx.Model(evdetails).Where("id = ?", evdetails[j].ID).Update("ticket_quota", deductedQuota).Error
						if err != nil {
							return err
						}
					}
					break
				}
			}
		}

		if err := tx.Create(&req).Error; err != nil {
			return err
		}

		return nil
	})

	repo.PopUserInOrderRoom(req.UserID, req.EventID, ctx)

	return err
}

func (repo *Implementation) isValidUniqueId(req models.BookingTicket, ctx context.Context) bool {
	orderRes := repo.GetUserInOrderRoom(req.UserID, req.EventID, ctx)
	if orderRes == "" {
		log.Println("User not found in Order Room")
		return false
	}
	datas := make(map[string]string)
	json.Unmarshal([]byte(orderRes), &datas)
	if datas["queueUniqueCode"] != req.QueueUniqueCode {
		log.Println("Unique code different")
		return false
	}
	return true
}

func (repo *Implementation) CheckBookingPeriod(ctx context.Context) {
	log.Println("Scheduler start")
	bookingRetention := 10
	var bookList []*models.BookingTicket
	joinQ := "LEFT JOIN purchased_tickets ON booking_tickets.id = purchased_tickets.booking_ticket_id"
	whereQ := "booking_status = ? AND purchased_tickets.id IS NULL AND CAST(EXTRACT(EPOCH FROM (NOW() - \"booking_tickets\".\"created_at\")) /60 as INTEGER) > ?"
	err := repo.db.WithContext(ctx).Preload("BookingDetails").Joins(joinQ).Where(whereQ, "active", bookingRetention).Find(&bookList).Error

	if err != nil {
		log.Println(err)
		return
	}
	retEventQuota := make(map[uint]uint)
	var eventIds []uint
	if bookList != nil {
		for i := 0; i < len(bookList); i++ {
			bookList[i].BookingStatus = "expired"

			for j := 0; j < len(bookList[i].BookingDetails); j++ {
				tempDetail := bookList[i].BookingDetails[j]

				val, found := retEventQuota[tempDetail.EventDetailID]
				if !found {
					val = 0
					eventIds = append(eventIds, tempDetail.EventDetailID)
				}
				retEventQuota[tempDetail.EventDetailID] = val + tempDetail.Qty
			}
		}
		eventDetails, err := repo.FindEventDetailsByIds(eventIds, context.TODO())
		if err != nil {
			return
		}

		for i := 0; i < len(eventDetails); i++ {
			eventDetails[i].TicketQuota += retEventQuota[eventDetails[i].ID]
		}

		if err2 := repo.db.Save(eventDetails).Error; err2 != nil {
			log.Println(err2)
			return
		}

		if err2 := repo.db.Save(bookList).Error; err2 != nil {
			log.Println(err2)
			return
		}
	}

}
