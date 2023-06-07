package repo

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
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

func TestImplementation_SaveBooking_shouldSuccess(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	implObj := NewImplementation(db, nil)
	os.Setenv("admin_fee", "2000")
	evRes := sqlmock.NewRows([]string{"id", "deleted_at"}).
		AddRow(1, nil)
	evDetRes := sqlmock.NewRows([]string{"id", "event_id", "ticket_quota", "deleted_at"}).
		AddRow(1, 1, 100, nil)

	//addRow := sqlmock.NewRows([]string{"id"}).AddRow("1")
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
	defer sqlDB.Close()

	implObj := NewImplementation(db, nil)
	os.Setenv("admin_fee", "2000")
	evRes := sqlmock.NewRows([]string{"id", "deleted_at"}).
		AddRow(1, nil)
	evDetRes := sqlmock.NewRows([]string{"id", "event_id", "ticket_quota", "deleted_at"}).
		AddRow(1, 1, 0, nil)

	//addRow := sqlmock.NewRows([]string{"id"}).AddRow("1")
	mock.ExpectBegin()
	mock.ExpectQuery(evSQL).WillReturnRows(evRes)
	mock.ExpectQuery(evDetailSQL).WillReturnRows(evDetRes)

	mock.ExpectRollback()

	var reqBook models.BookingTicket
	reqBook.EventID = 1
	reqBook.BookingDetails = genBookDetail()
	err := implObj.SaveBooking(reqBook, context.TODO())
	assert.Equal(t, err.Error(), "ticket quota not enough")
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
