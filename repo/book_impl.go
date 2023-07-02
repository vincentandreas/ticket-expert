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
	err := tx.WithContext(ctx).Model(evdetails).Where("id = ?", id).Update("ticket_remaining", quota).Error
	return err
}

var validateQUniqueId = IsValidUniqueId
var popUserHelper = PopUserInOrderRoom

func (repo *Implementation) SaveBooking(req models.BookingTicket, ctx context.Context) error {
	var grandTotal float64 = 0
	bookDetails := req.BookingDetails

	if !validateQUniqueId(repo, req, ctx) {
		return errors.New("Failed, Queue Unique Id Not Match")
	}
	var evDetailIds []uint
	//getting the event detail list
	for i := 0; i < len(bookDetails); i++ {
		evDetailIds = append(evDetailIds, bookDetails[i].EventDetailID)
	}

	priceMap, err := repo.FindEvDetailPrice(evDetailIds, ctx)
	if err != nil {
		return err
	}

	for i := 0; i < len(bookDetails); i++ {
		ticketPrice, ok := priceMap[bookDetails[i].EventDetailID]
		// If the key exists
		if !ok {
			return errors.New("event detail id not found")
		}

		price, _ := strconv.ParseFloat(ticketPrice, 64)

		total := price * float64(bookDetails[i].Qty)
		grandTotal += total
		bookDetails[i].Price = ticketPrice
		bookDetails[i].SubTotal = strconv.FormatFloat(total, 'f', -1, 64)
	}
	admEnv := os.Getenv("admin_fee")
	admFee, _ := strconv.ParseFloat(admEnv, 64)
	grandTotal += admFee
	req.TotalPrice = strconv.FormatFloat(grandTotal, 'f', -1, 64)
	req.AdminFee = admEnv

	err = repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var evData models.Event
		tx.Preload("EventDetails").Where("events.id = ?", req.EventID).Find(&evData)

		evdetails := evData.EventDetails

		for i := 0; i < len(bookDetails); i++ {
			for j := 0; j < len(evdetails); j++ {
				if evdetails[j].ID == bookDetails[i].EventDetailID {
					if evdetails[j].TicketRemaining < bookDetails[i].Qty {
						return errors.New("ticket remaining is not enough")
					} else {
						deductedQuota := evdetails[j].TicketRemaining - bookDetails[i].Qty
						err := repo.UpdTicketQty(evdetails[j].ID, deductedQuota, tx, ctx)
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

	popUserHelper(repo, req.UserID, req.EventID, ctx)

	return err
}

var helperGetUser = GetUserInOrderRoom

func IsValidUniqueId(repo *Implementation, req models.BookingTicket, ctx context.Context) bool {
	orderRes := helperGetUser(repo, req.UserID, req.EventID, ctx)
	if orderRes == "" {
		log.Println("User not found in Order Room")
		return false
	}
	datas := make(map[string]string)
	json.Unmarshal([]byte(orderRes), &datas)
	if datas["qUniqueCode"] != req.QUniqueCode {
		log.Println("Unique code different")
		return false
	}
	return true
}

var helperGetBooklistExceed = GetBookingExceedRetention
var helperFindEvDetails = FindEventDetailsByIds

func (repo *Implementation) CheckBookingPeriodically(ctx context.Context) {
	log.Println("Scheduler start")
	bookingRetention := 10
	bookList, err := helperGetBooklistExceed(ctx, repo, bookingRetention)

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
		eventDetails, err := helperFindEvDetails(repo, eventIds, context.TODO())
		if err != nil {
			return
		}

		for i := 0; i < len(eventDetails); i++ {
			eventDetails[i].TicketRemaining += retEventQuota[eventDetails[i].ID]
		}

		err = repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

			if err2 := tx.Save(eventDetails).Error; err2 != nil {
				log.Println(err2)
				return err2
			}

			if err2 := tx.Save(bookList).Error; err2 != nil {
				log.Println(err2)
				return err2
			}
			return nil
		})
		log.Println(err)
	}

}

func (repo *Implementation) FindBookingByUserId(userId uint, ctx context.Context) ([]*models.ShowBooking, error) {
	var showBooks []*models.ShowBooking
	selSql := "events.event_name, booking_tickets.q_unique_code, booking_tickets.booking_status, booking_tickets.total_price"
	joinSql := "JOIN events ON booking_tickets.event_id = events.id"
	err := repo.db.WithContext(ctx).Table("booking_tickets").Select(selSql).Joins(joinSql).Where("booking_tickets.user_id = ?", userId).Scan(&showBooks).Error
	return showBooks, err
}

func GetBookingExceedRetention(ctx context.Context, repo *Implementation, bookingRetention int) ([]*models.BookingTicket, error) {
	var bookList []*models.BookingTicket
	joinQ := "LEFT JOIN purchased_tickets ON booking_tickets.id = purchased_tickets.booking_ticket_id"
	whereQ := "booking_status = ? AND purchased_tickets.id IS NULL AND CAST(EXTRACT(EPOCH FROM (NOW() - \"booking_tickets\".\"created_at\")) /60 as INTEGER) > ?"
	err := repo.db.WithContext(ctx).Preload("BookingDetails").Joins(joinQ).Where(whereQ, "active", bookingRetention).Find(&bookList).Error
	return bookList, err
}

func (repo *Implementation) GetBookingByUniqCode(ctx context.Context, uniqCode string) (*models.BookingTicket, error) {
	var book *models.BookingTicket
	err := repo.db.WithContext(ctx).Where("q_unique_code = ?", uniqCode).First(&book).Error
	return book, err
}

func (repo *Implementation) GetBookingDataByUniqCode(ctx context.Context, uniqCode string) (models.PurchaseDetails, error) {
	var purchaseDetails models.PurchaseDetails
	var bookingTicket models.BookingTicket
	var event models.Event
	var eventDetails []models.EventDetail

	if err := repo.db.WithContext(ctx).Preload("PurchasedTicket").Preload("BookingDetails").Where("q_unique_code = ?", uniqCode).First(&bookingTicket).Error; err != nil {
		return purchaseDetails, err
	}

	if err := repo.db.WithContext(ctx).First(&event, bookingTicket.EventID).Error; err != nil {
		return purchaseDetails, err
	}

	var evDetIds []uint
	var tdetails []models.TicketDetails

	for i := 0; i < len(bookingTicket.BookingDetails); i++ {
		evDetIds = append(evDetIds, bookingTicket.BookingDetails[i].EventDetailID)
	}

	if err := repo.db.WithContext(ctx).Where("id in ?", evDetIds).Find(&eventDetails).Error; err != nil {
		return purchaseDetails, err
	}

	for i := 0; i < len(bookingTicket.BookingDetails); i++ {
		for j := 0; j < len(eventDetails); j++ {
			bd := bookingTicket.BookingDetails[i]
			ed := eventDetails[j]
			if bd.EventDetailID == ed.ID {
				tmp := models.TicketDetails{
					Qty:         bd.Qty,
					Price:       bd.Price,
					SubTotal:    bd.SubTotal,
					TicketClass: ed.TicketClass,
				}
				tdetails = append(tdetails, tmp)
				break
			}
		}
	}

	purchaseDetails.AdminFee = bookingTicket.AdminFee
	purchaseDetails.EventCategory = event.EventCategory
	purchaseDetails.EventName = event.EventName
	purchaseDetails.QUniqueCode = bookingTicket.QUniqueCode
	purchaseDetails.BookingStatus = bookingTicket.BookingStatus
	purchaseDetails.TotalPrice = bookingTicket.TotalPrice
	purchaseDetails.TicketDetails = tdetails
	return purchaseDetails, nil
}
