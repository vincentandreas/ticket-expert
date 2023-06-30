package repo

import (
	"context"
	"encoding/json"
	"errors"
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

func TestCountTotalPeopleInWaitingRoom_shouldNotFound(t *testing.T) {
	db, mock := redismock.NewClientMock()
	implObj := NewImplementation(nil, db)
	mock.Regexp().ExpectLLen("QEvent2").RedisNil()
	wrTotal := CountTotalPeopleInWaitingRoom(implObj, 2, context.TODO())
	assert.Equal(t, wrTotal, int64(-1))
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCountTotalPeopleInWaitingRoom_shouldFound(t *testing.T) {
	db, mock := redismock.NewClientMock()
	implObj := NewImplementation(nil, db)
	mock.Regexp().ExpectLLen("QEvent2").SetVal(3)
	wrTotal := CountTotalPeopleInWaitingRoom(implObj, 2, context.TODO())
	assert.Equal(t, wrTotal, int64(3))
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestImplementation_PopWaitingQueue_shouldNull(t *testing.T) {
	db, mock := redismock.NewClientMock()
	implObj := NewImplementation(nil, db)
	mock.Regexp().ExpectRPop("QEvent1").RedisNil()
	res := PopWaitingQueue(implObj, 1, context.TODO())
	assert.Equal(t, res, "")
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestImplementation_PopWaitingQueue_shouldFound(t *testing.T) {
	db, mock := redismock.NewClientMock()
	implObj := NewImplementation(nil, db)
	mock.Regexp().ExpectRPop("QEvent1").SetVal("abc123")
	res := PopWaitingQueue(implObj, 1, context.TODO())
	assert.Equal(t, res, "abc123")
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestImplementation_SaveUserInOrderRoom(t *testing.T) {
	db, mock := redismock.NewClientMock()
	implObj := NewImplementation(nil, db)
	mock.Regexp().ExpectHSet("Event1", "1", ".+").RedisNil()

	SaveUserInOrderRoom(implObj, 1, "1", "uniquniq", context.TODO())

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestPopUserInOrderRoom_shouldTrue(t *testing.T) {
	db, mock := redismock.NewClientMock()
	implObj := NewImplementation(nil, db)

	mock.Regexp().ExpectHDel("ORmEvent1", "1").SetVal(1)
	popres := PopUserInOrderRoom(implObj, 1, 1, context.TODO())
	assert.True(t, popres)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestPopUserInOrderRoom_shouldFalse(t *testing.T) {
	db, mock := redismock.NewClientMock()
	implObj := NewImplementation(nil, db)

	mock.Regexp().ExpectHDel("ORmEvent1", "1").SetErr(errors.New("fail"))
	popres := PopUserInOrderRoom(implObj, 1, 1, context.TODO())
	assert.False(t, popres)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetUserInOrderRoom_shouldFound(t *testing.T) {
	db, mock := redismock.NewClientMock()
	implObj := NewImplementation(nil, db)

	mck := make(map[string]string)
	mck["qUniqueCode"] = "123a"
	mckStr, _ := json.Marshal(mck)
	mock.Regexp().ExpectHGet("ORmEvent1", "1").SetVal(string(mckStr))

	res := GetUserInOrderRoom(implObj, 1, 1, context.TODO())
	assert.Equal(t, res, string(mckStr))
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetUserInOrderRoom_shouldNotFound(t *testing.T) {
	db, mock := redismock.NewClientMock()
	implObj := NewImplementation(nil, db)

	mock.Regexp().ExpectHGet("ORmEvent1", "1").SetErr(errors.New("failed"))

	res := GetUserInOrderRoom(implObj, 1, 1, context.TODO())
	assert.Equal(t, res, "")
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCountPeopleInOrderRoom_shouldSuccess(t *testing.T) {
	db, mock := redismock.NewClientMock()
	implObj := NewImplementation(nil, db)

	mock.Regexp().ExpectHLen("ORmEvent1").SetVal(7)
	popRes := CountPeopleInOrderRoom(implObj, 1, context.TODO())

	assert.Equal(t, popRes, int64(7))
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCountPeopleInOrderRoom_shouldError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	implObj := NewImplementation(nil, db)

	mock.Regexp().ExpectHLen("ORmEvent1").SetErr(errors.New("failed"))
	popRes := CountPeopleInOrderRoom(implObj, 1, context.TODO())

	assert.Equal(t, popRes, int64(-1))
	assert.Nil(t, mock.ExpectationsWereMet())
}
