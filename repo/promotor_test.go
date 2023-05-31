package repo

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"ticket-expert/models"
)

func TestAddPromotor(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	implObj := NewImplementation(db)

	expectedSQL := "INSERT INTO \"promotors\" (.+) VALUES (.+)"
	mock.ExpectBegin()
	mock.ExpectQuery(expectedSQL).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	mock.ExpectCommit()
	var reqPromotor models.Promotor
	implObj.SavePromotor(reqPromotor, context.TODO())
	assert.Nil(t, mock.ExpectationsWereMet())
}
