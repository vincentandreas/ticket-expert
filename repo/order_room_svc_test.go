package repo

import (
	"context"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImplementation_CheckOrderRoom_shouldRetEmptyArr(t *testing.T) {
	helperCountPeople = func(repo *Implementation, eventId uint, ctx context.Context) int64 {
		return 1
	}
	helperPopQueue = func(repo *Implementation, eventId uint, ctx context.Context) string {
		return ""
	}
	db, _ := redismock.NewClientMock()
	implObj := NewImplementation(nil, db)

	qCodes := implObj.CheckOrderRoom(1, context.TODO())
	assert.Equal(t, len(qCodes), 0)
}

func TestImplementation_CheckOrderRoom_shouldRetQCodes(t *testing.T) {
	ctr := 0
	helperCountPeople = func(repo *Implementation, eventId uint, ctx context.Context) int64 {
		ctr += 1
		return int64(ctr)
	}
	helperPopQueue = func(repo *Implementation, eventId uint, ctx context.Context) string {
		return "{\"user_id\":1,\"q_unique_code\":\"123123123\",\"event_id\":1}"
	}
	db, _ := redismock.NewClientMock()
	implObj := NewImplementation(nil, db)

	qCodes := implObj.CheckOrderRoom(1, context.TODO())
	assert.Equal(t, len(qCodes), 1)
}
