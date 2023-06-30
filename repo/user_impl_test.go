package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"ticket-expert/models"
	"time"
)

func DbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		t.Fatal(err)
	}
	return sqldb, gormdb, mock
}

func TestFindUser(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	implObj := NewImplementation(db, nil)
	timenow := time.Time{}
	users := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "full_name", "user_name", "password"}).
		AddRow(1, timenow, timenow, timenow, "user", "user", "passwd")

	expectedSQL := "SELECT (.+) FROM \"users\" WHERE id =(.+)"
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)
	_, res := implObj.FindUserById(1, context.TODO())
	assert.Nil(t, res)

	_, res2 := implObj.FindUserById(2, context.TODO())
	assert.True(t, errors.Is(res2, gorm.ErrRecordNotFound))
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestAddUser(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	implObj := NewImplementation(db, nil)

	expectedSQL := "INSERT INTO \"users\" (.+) VALUES (.+)"
	mock.ExpectBegin()
	addRow := sqlmock.NewRows([]string{"id"}).AddRow("1")
	mock.ExpectQuery(expectedSQL).WillReturnRows(addRow)
	mock.ExpectCommit()
	var reqUser models.User
	implObj.SaveUser(reqUser, context.TODO())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestImplementation_Login_shouldScs(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	implObj := NewImplementation(db, nil)
	timenow := time.Time{}

	userRes := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "full_name", "user_name", "password"}).
		AddRow(1, timenow, timenow, timenow, "user", "user", "passwd")

	expectedSQL := "SELECT (.+) FROM \"users\" WHERE .+"
	mock.ExpectQuery(expectedSQL).WillReturnRows(userRes)

	userObj := models.UserLogin{
		UserName: "user",
		Password: "passwd",
	}
	_, err := implObj.Login(userObj, context.TODO())
	assert.Nil(t, err)
}

func TestImplementation_Login_shouldFalse(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	implObj := NewImplementation(db, nil)

	userRes := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "full_name", "user_name", "password"})

	expectedSQL := "SELECT (.+) FROM \"users\" WHERE .+"
	mock.ExpectQuery(expectedSQL).WillReturnRows(userRes)

	userObj := models.UserLogin{
		UserName: "user",
		Password: "passwd",
	}
	_, err := implObj.Login(userObj, context.TODO())
	assert.Equal(t, err.Error(), "record not found")
}
