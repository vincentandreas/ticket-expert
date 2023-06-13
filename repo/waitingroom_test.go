package repo

import (
	"context"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"ticket-expert/models"
)

func TestImplementation_SaveWaitingQueue(t *testing.T) {
	db, mock := redismock.NewClientMock()
	implObj := NewImplementation(nil, db)
	var wuser models.NewWaitingUser
	mock.Regexp().ExpectLPush("QEvent1", `.+`).RedisNil()
	wuser.EventId = 1
	implObj.SaveWaitingQueue(wuser, context.TODO())

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestImplementation_PopWaitingQueue(t *testing.T) {
	db, mock := redismock.NewClientMock()
	implObj := NewImplementation(nil, db)
	mock.Regexp().ExpectRPop("QEvent1").RedisNil()
	PopWaitingQueue(implObj, 1, context.TODO())

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestImplementation_SaveUserInOrderRoom(t *testing.T) {
	db, mock := redismock.NewClientMock()
	implObj := NewImplementation(nil, db)
	mock.Regexp().ExpectHSet("Event1", "1", ".+").RedisNil()

	SaveUserInOrderRoom(implObj, 1, "1", "uniquniq", context.TODO())

	assert.Nil(t, mock.ExpectationsWereMet())
}
