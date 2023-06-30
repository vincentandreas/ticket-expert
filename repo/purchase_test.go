package repo

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"ticket-expert/models"
)

func TestImplementation_SavePurchase_shouldSuccess(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	implObj := NewImplementation(db, nil)
	purchaseSql := "SELECT .+ FROM \"purchased_tickets\""
	bookingSql := "SELECT .+ FROM \"booking_tickets\""
	insertPurchaseSQL := "INSERT INTO \"purchased_tickets\""
	updateBookSql := "UPDATE \"booking_tickets\" SET .+ WHERE .+"

	purchaseSel := sqlmock.NewRows([]string{"id", "deleted_at"})
	bookingSel := sqlmock.NewRows([]string{"id", "booking_status"}).AddRow(1, "active")
	mock.ExpectBegin()
	mock.ExpectQuery(bookingSql).WillReturnRows(bookingSel)
	mock.ExpectQuery(purchaseSql).WillReturnRows(purchaseSel)
	mock.ExpectQuery(insertPurchaseSQL).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	mock.ExpectExec(updateBookSql).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	var reqPur models.PurchasedTicket

	err := implObj.SavePurchase(reqPur, context.TODO())
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestImplementation_SavePurchase_shouldFailed_whenBookingExp(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	implObj := NewImplementation(db, nil)
	purchaseSql := "SELECT .+ FROM \"purchased_tickets\""
	bookingSql := "SELECT .+ FROM \"booking_tickets\""
	//insertPurchaseSQL := "INSERT INTO \"purchased_tickets\""

	purchaseSel := sqlmock.NewRows([]string{"id", "deleted_at"})
	bookingSel := sqlmock.NewRows([]string{"id", "booking_status"}).AddRow(1, "expired")
	mock.ExpectBegin()
	mock.ExpectQuery(bookingSql).WillReturnRows(bookingSel)
	mock.ExpectQuery(purchaseSql).WillReturnRows(purchaseSel)
	//mock.ExpectQuery(insertPurchaseSQL).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	mock.ExpectRollback()
	var reqPur models.PurchasedTicket

	err := implObj.SavePurchase(reqPur, context.TODO())
	assert.Equal(t, err.Error(), "booking status not active")
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestImplementation_SavePurchase_shouldFailed_whenAlreadyPurchased(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	implObj := NewImplementation(db, nil)
	purchaseSql := "SELECT .+ FROM \"purchased_tickets\""
	bookingSql := "SELECT .+ FROM \"booking_tickets\""
	//insertPurchaseSQL := "INSERT INTO \"purchased_tickets\""

	purchaseSel := sqlmock.NewRows([]string{"id", "booking_ticket_id", "payment_status"}).AddRow(1, 1, "success")
	bookingSel := sqlmock.NewRows([]string{"id", "booking_status"}).AddRow(1, "active")
	mock.ExpectBegin()
	mock.ExpectQuery(bookingSql).WillReturnRows(bookingSel)
	mock.ExpectQuery(purchaseSql).WillReturnRows(purchaseSel)

	mock.ExpectRollback()
	var reqPur models.PurchasedTicket

	err := implObj.SavePurchase(reqPur, context.TODO())
	assert.Equal(t, err.Error(), "this booking already purchased")
	assert.Nil(t, mock.ExpectationsWereMet())
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
