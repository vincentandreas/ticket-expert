package repo

import (
	"context"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
	"ticket-expert/models"
)

var updEventDetailSQL = "UPDATE \"event_details\" SET .+"
var evSQL = "SELECT .+ FROM \"events\" WHERE .+"
var evDetailSQL = "SELECT .+ FROM \"event_details\" WHERE .+"
var insertBookSQL = "INSERT INTO \"booking_tickets\""
var insertBookDetSQL = "INSERT INTO \"booking_details\""

func TestImplementation_UpdTicketQty(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	implObj := NewImplementation(db, nil)

	//kk := "UPDATE \"event_details\" SET .+"
	mock.ExpectBegin()
	mock.ExpectExec(updEventDetailSQL).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := implObj.UpdTicketQty(1, 99, db, context.TODO())
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestImplementation_SaveBooking_failed_when_quid_notvalid(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer func() {
		sqlDB.Close()
		validateQUniqueId = IsValidUniqueId
	}()

	validateQUniqueId = func(repo *Implementation, req models.BookingTicket, ctx context.Context) bool {
		return false
	}

	implObj := NewImplementation(db, nil)

	var reqBook models.BookingTicket
	reqBook.EventID = 1
	reqBook.BookingDetails = genBookDetail()
	err := implObj.SaveBooking(reqBook, context.TODO())
	assert.Equal(t, err.Error(), "Failed, Queue Unique Id Not Match")
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestImplementation_SaveBooking_shouldSuccess(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer func() {
		sqlDB.Close()
		validateQUniqueId = IsValidUniqueId
		popUserHelper = PopUserInOrderRoom
	}()

	//mock validate uid
	validateQUniqueId = func(repo *Implementation, req models.BookingTicket, ctx context.Context) bool {
		return true
	}

	popUserHelper = func(repo *Implementation, userId uint, eventId uint, ctx context.Context) bool {
		log.Println("Mocking pop helper")
		return true
	}

	implObj := NewImplementation(db, nil)
	os.Setenv("admin_fee", "2000")
	evRes := sqlmock.NewRows([]string{"id", "deleted_at"}).
		AddRow(1, nil)
	evDetRes := sqlmock.NewRows([]string{"id", "event_id", "ticket_quota", "ticket_remaining", "deleted_at"}).
		AddRow(1, 1, 100, 100, nil)
	checkPriceRes := sqlmock.NewRows([]string{"id", "ticket_price"}).
		AddRow(1, "10000")
	//addRow := sqlmock.NewRows([]string{"id"}).AddRow("1")

	checkPriceSql := "SELECT .+ FROM \"event_details\" WHERE .+"

	mock.ExpectQuery(checkPriceSql).WillReturnRows(checkPriceRes)

	mock.ExpectBegin()
	mock.ExpectQuery(evSQL).WillReturnRows(evRes)
	mock.ExpectQuery(evDetailSQL).WillReturnRows(evDetRes)
	mock.ExpectExec(updEventDetailSQL).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(insertBookSQL).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	mock.ExpectQuery(insertBookDetSQL).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	mock.ExpectCommit()

	var reqBook models.BookingTicket
	reqBook.EventID = 1
	reqBook.BookingDetails = genBookDetail()
	err := implObj.SaveBooking(reqBook, context.TODO())
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestImplementation_SaveBooking_shouldFail_whenQuotaNotEnough(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer func() {
		sqlDB.Close()
		validateQUniqueId = IsValidUniqueId
		popUserHelper = PopUserInOrderRoom

	}()

	validateQUniqueId = func(repo *Implementation, req models.BookingTicket, ctx context.Context) bool {
		return true
	}
	popUserHelper = func(repo *Implementation, userId uint, eventId uint, ctx context.Context) bool {
		log.Println("Mocking pop helper")
		return true
	}

	implObj := NewImplementation(db, nil)
	os.Setenv("admin_fee", "2000")
	evRes := sqlmock.NewRows([]string{"id", "deleted_at"}).
		AddRow(1, nil)
	evDetRes := sqlmock.NewRows([]string{"id", "event_id", "ticket_quota", "deleted_at"}).
		AddRow(1, 1, 0, nil)
	checkPriceRes := sqlmock.NewRows([]string{"id", "ticket_price"}).
		AddRow(1, "10000")
	checkPriceSql := "SELECT .+ FROM \"event_details\" WHERE .+"

	//addRow := sqlmock.NewRows([]string{"id"}).AddRow("1")
	mock.ExpectQuery(checkPriceSql).WillReturnRows(checkPriceRes)
	mock.ExpectBegin()
	mock.ExpectQuery(evSQL).WillReturnRows(evRes)
	mock.ExpectQuery(evDetailSQL).WillReturnRows(evDetRes)

	mock.ExpectRollback()

	var reqBook models.BookingTicket
	reqBook.EventID = 1
	reqBook.BookingDetails = genBookDetail()
	err := implObj.SaveBooking(reqBook, context.TODO())
	assert.Equal(t, err.Error(), "ticket remaining is not enough")
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestImplementation_getBookingExceedRetention(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer func() {
		sqlDB.Close()
	}()
	implObj := NewImplementation(db, nil)

	res := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "total_price", "admin_fee", "user_id", "booking_status", "q_unique_code", "event_id"})

	mock.ExpectQuery("SELECT .+ FROM \"booking_tickets\" LEFT JOIN purchased_tickets ON booking_tickets.id = purchased_tickets.booking_ticket_id WHERE .+").WillReturnRows(res)

	GetBookingExceedRetention(context.TODO(), implObj, 10)
	assert.Nil(t, mock.ExpectationsWereMet())

}

func TestImplementation_CheckBookingPeriod(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer func() {
		sqlDB.Close()
		helperGetBooklistExceed = GetBookingExceedRetention
	}()

	helperGetBooklistExceed = func(ctx context.Context, repo *Implementation, bookingRetention int64) ([]*models.BookingTicket, error) {
		var booklist []*models.BookingTicket
		sgl := models.BookingTicket{
			BookingStatus:  "active",
			BookingDetails: genBookDetail(),
		}
		booklist = append(booklist, &sgl)
		return booklist, nil
	}
	helperFindEvDetails = func(repo *Implementation, ids []uint, ctx context.Context) ([]*models.EventDetail, error) {
		var eds []*models.EventDetail
		sgl := models.EventDetail{
			Model: gorm.Model{
				ID: 1,
			},
			TicketQuota:     100,
			TicketRemaining: 57,
		}
		eds = append(eds, &sgl)
		return eds, nil
	}

	res := sqlmock.NewRows([]string{""})

	implObj := NewImplementation(db, nil)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO \"event_details\" .+ VALUES .+ ON CONFLICT (.+) DO UPDATE SET.+").WillReturnRows(res)
	mock.ExpectQuery("INSERT INTO \"booking_tickets\" .+ VALUES .+ ON CONFLICT (.+) DO UPDATE SET.+").WillReturnRows(res)
	mock.ExpectQuery("INSERT INTO \"booking_details\" .+ VALUES .+ ON CONFLICT (.+) DO UPDATE SET.+").WillReturnRows(res)
	mock.ExpectCommit()
	implObj.CheckBookingPeriod(context.TODO())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestImplementation_FindBookingByUserId(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer func() {
		sqlDB.Close()
	}()
	implObj := NewImplementation(db, nil)

	res := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "total_price", "admin_fee", "user_id", "booking_status", "q_unique_code", "event_id"})

	mock.ExpectQuery("SELECT .+ FROM \"booking_tickets\" JOIN events ON booking_tickets.event_id = events.id WHERE .+").WillReturnRows(res)
	implObj.FindBookingByUserId(1, context.TODO())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestIsValidUniqueId_shouldTrue(t *testing.T) {
	sqlDB, db, _ := DbMock(t)
	defer func() {
		sqlDB.Close()
	}()
	implObj := NewImplementation(db, nil)

	helperGetUser = func(repo *Implementation, userId uint, eventId uint, ctx context.Context) string {
		log.Println("Dummy get user")
		mck := make(map[string]string)
		mck["qUniqueCode"] = "123a"
		mckStr, err := json.Marshal(mck)
		if err != nil {
			return ""
		}
		return string(mckStr)
	}
	bookTick := models.BookingTicket{
		QUniqueCode: "123a",
	}
	res := IsValidUniqueId(implObj, bookTick, context.TODO())
	assert.True(t, res)
}

func TestIsValidUniqueId_shouldFalse_whenQDiff(t *testing.T) {
	sqlDB, db, _ := DbMock(t)
	defer func() {
		sqlDB.Close()
	}()
	implObj := NewImplementation(db, nil)

	helperGetUser = func(repo *Implementation, userId uint, eventId uint, ctx context.Context) string {
		log.Println("Dummy get user")
		mck := make(map[string]string)
		mck["qUniqueCode"] = "123a"
		mckStr, err := json.Marshal(mck)
		if err != nil {
			return ""
		}
		return string(mckStr)
	}
	bookTick := models.BookingTicket{
		QUniqueCode: "888",
	}
	res := IsValidUniqueId(implObj, bookTick, context.TODO())
	assert.False(t, res)
}

func TestIsValidUniqueId_shouldNotFound(t *testing.T) {
	sqlDB, db, _ := DbMock(t)
	defer func() {
		sqlDB.Close()
	}()
	implObj := NewImplementation(db, nil)

	helperGetUser = func(repo *Implementation, userId uint, eventId uint, ctx context.Context) string {
		return ""
	}
	bookTick := models.BookingTicket{
		QUniqueCode: "123a",
	}
	res := IsValidUniqueId(implObj, bookTick, context.TODO())
	assert.False(t, res)
}
func TestImplementation_GetBookingByUniqCode_shouldFound(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	implObj := NewImplementation(db, nil)
	selBooksql := "SELECT .+ FROM \"booking_tickets\" WHERE q_unique_code = .+"
	bookRes := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(selBooksql).WillReturnRows(bookRes)

	res, err := implObj.GetBookingByUniqCode(context.TODO(), "")
	assert.Nil(t, err)
	assert.Equal(t, res.ID, uint(1))
}

func TestImplementation_GetBookingByUniqCode_shouldNotFound(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	implObj := NewImplementation(db, nil)
	selBooksql := "SELECT .+ FROM \"booking_tickets\" WHERE q_unique_code = .+"
	bookRes := sqlmock.NewRows([]string{"id"})
	mock.ExpectQuery(selBooksql).WillReturnRows(bookRes)

	_, err := implObj.GetBookingByUniqCode(context.TODO(), "")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "record not found")
}

func TestImplementation_GetBookingDataByUniqCode_shouldHandle_notFound(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	implObj := NewImplementation(db, nil)

	purchaseSql := "SELECT .+ FROM \"purchased_tickets\" WHERE .+"
	bookingSql := "SELECT .+ FROM \"booking_tickets\" WHERE .+"
	bookDetSql := "SELECT .+ FROM \"booking_details\" WHERE .+"
	evSql := "SELECT .+ FROM \"events\" WHERE .+"

	sqlRes := sqlmock.NewRows([]string{"id"}).AddRow(1)
	bookDetRes := sqlmock.NewRows([]string{"id", "event_detail_id", "booking_ticket_id"}).AddRow(1, 3, 1)
	evRes := sqlmock.NewRows([]string{"id", "event_detail_id"})

	mock.ExpectQuery(bookingSql).WillReturnRows(sqlRes)
	mock.ExpectQuery(bookDetSql).WillReturnRows(bookDetRes)
	mock.ExpectQuery(purchaseSql).WillReturnRows(sqlRes)
	mock.ExpectQuery(evSql).WillReturnRows(evRes)

	_, err := implObj.GetBookingDataByUniqCode(context.TODO(), "")
	assert.Equal(t, err.Error(), "record not found")
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestImplementation_GetBookingDataByUniqCode(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	implObj := NewImplementation(db, nil)

	purchaseSql := "SELECT .+ FROM \"purchased_tickets\" WHERE .+"
	bookingSql := "SELECT .+ FROM \"booking_tickets\" WHERE .+"
	bookDetSql := "SELECT .+ FROM \"booking_details\" WHERE .+"
	evSql := "SELECT .+ FROM \"events\" WHERE .+"
	evDetSql := "SELECT .+ FROM \"event_details\" WHERE .+"

	sqlRes := sqlmock.NewRows([]string{"id"}).AddRow(1)
	bookDetRes := sqlmock.NewRows([]string{"id", "event_detail_id", "booking_ticket_id"}).AddRow(1, 3, 1)
	evRes := sqlmock.NewRows([]string{"id", "event_detail_id"}).AddRow(1, 3)
	evDetRes := sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(3)
	mock.ExpectQuery(bookingSql).WillReturnRows(sqlRes)
	mock.ExpectQuery(bookDetSql).WillReturnRows(bookDetRes)
	mock.ExpectQuery(purchaseSql).WillReturnRows(sqlRes)
	mock.ExpectQuery(evSql).WillReturnRows(evRes)
	mock.ExpectQuery(evDetSql).WillReturnRows(evDetRes)

	purchaseData, err := implObj.GetBookingDataByUniqCode(context.TODO(), "")
	assert.Nil(t, err)
	assert.Equal(t, len(purchaseData.TicketDetails), 1)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func genBookDetail() []*models.BookingDetail {
	var arrdetail []*models.BookingDetail

	var reqbook models.BookingDetail
	reqbook.ID = 1
	reqbook.Qty = 3
	reqbook.Price = "1000"
	reqbook.EventDetailID = 1
	arrdetail = append(arrdetail, &reqbook)
	return arrdetail
}
