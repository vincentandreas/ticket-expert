package repo

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"ticket-expert/models"
)

func TestImplementation_FindEventByCondition(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	implObj := NewImplementation(db)
	evRes := sqlmock.NewRows([]string{"id", "dummy_detail"}).
		AddRow(1, "user")

	evSQL := "SELECT .+ FROM \"events\" LEFT JOIN promotors on events.promotor_id = promotors.id WHERE (.+)"
	mock.ExpectQuery(evSQL).WillReturnRows(evRes)
	_, err := implObj.FindEventByCondition("jakarta", "music", context.TODO())
	assert.Nil(t, err)
}

func TestImplementation_FindEventById(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	implObj := NewImplementation(db)
	evRes := sqlmock.NewRows([]string{"id", "dummy_detail"}).
		AddRow(1, "user")

	evSQL := "SELECT .+ FROM \"events\" WHERE .+"
	evDetailSQL := "SELECT .+ FROM \"event_details\" WHERE .+"
	mock.ExpectQuery(evSQL).WillReturnRows(evRes)
	mock.ExpectQuery(evDetailSQL).WillReturnRows(evRes)
	_, err := implObj.FindByEventId("1", context.TODO())
	assert.Nil(t, err)
}

func TestImplementation_SaveEvent(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	implObj := NewImplementation(db)

	eventSql := "INSERT INTO \"events\" (.+) VALUES (.+)"
	eventDtlSql := "INSERT INTO \"event_details\" (.+) VALUES (.+)"
	mock.ExpectBegin()
	mock.ExpectQuery(eventSql).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	mock.ExpectQuery(eventDtlSql).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	mock.ExpectCommit()
	var reqEvent models.Event

	var reqEvDetail models.EventDetail
	reqEvDetail.TicketQuota = 100
	reqEvDetail.TicketPrice = "3000"

	var arrEvdetail []*models.EventDetail
	arrEvdetail = append(arrEvdetail, &reqEvDetail)
	reqEvent.EventDetails = arrEvdetail
	implObj.SaveEvent(reqEvent, context.TODO())
	assert.Nil(t, mock.ExpectationsWereMet())
}
