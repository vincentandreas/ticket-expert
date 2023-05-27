package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

func dbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return sqldb, gormdb, mock
}

func TestAddUser(t *testing.T) {

	sqlDB, db, mock := dbMock(t)
	defer sqlDB.Close()

	implObj := NewImplementation(db)
	timenow := time.Time{}
	users := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "full_name", "user_name", "password"}).
		AddRow(1, timenow, timenow, timenow, "user", "user", "passwd")

	expectedSQL := "SELECT (.+) FROM \"users\" WHERE id =(.+)"
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)
	_, res := implObj.FindUserById(1, context.TODO())
	assert.Nil(t, res.Error)

	_, res2 := implObj.FindUserById(2, context.TODO())
	assert.True(t, errors.Is(res2.Error, gorm.ErrRecordNotFound))
	assert.Nil(t, mock.ExpectationsWereMet())

}
