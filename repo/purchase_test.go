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

	implObj := NewImplementation(db)
	purchaseSql := "SELECT .+ FROM \"purchased_tickets\""
	bookingSql := "SELECT .+ FROM \"booking_tickets\""
	insertPurchaseSQL := "INSERT INTO \"purchased_tickets\""

	purchaseSel := sqlmock.NewRows([]string{"id", "deleted_at"})
	bookingSel := sqlmock.NewRows([]string{"id", "booking_status"}).AddRow(1, "active")
	mock.ExpectBegin()
	mock.ExpectQuery(bookingSql).WillReturnRows(bookingSel)
	mock.ExpectQuery(purchaseSql).WillReturnRows(purchaseSel)
	mock.ExpectQuery(insertPurchaseSQL).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	mock.ExpectCommit()
	var reqPur models.PurchasedTicket

	err := implObj.SavePurchase(reqPur, context.TODO())
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestImplementation_SavePurchase_shouldFailed_whenBookingExp(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	implObj := NewImplementation(db)
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

	implObj := NewImplementation(db)
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
