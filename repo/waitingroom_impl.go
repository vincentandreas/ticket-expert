package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"ticket-expert/models"
	"time"
)

func (repo *Implementation) SaveWaitingQueue(wuser models.NewWaitingUser, ctx context.Context) {
	key := genQueueEventStr(wuser.EventId)
	wuserJson, _ := json.Marshal(wuser)
	_, err := repo.redis.LPush(ctx, key, wuserJson).Result()
	if err != nil {
		log.Println(err)
		return
	}
}

func (repo *Implementation) CountTotalPeopleInWaitingRoom(eventId uint, ctx context.Context) int64 {
	key := genQueueEventStr(eventId)
	result, err := repo.redis.LLen(ctx, key).Result()
	if err != nil {
		return -1
	}
	return result
}

func PopWaitingQueue(repo *Implementation, eventId uint, ctx context.Context) string {
	key := genQueueEventStr(eventId)
	result, err := repo.redis.RPop(ctx, key).Result()
	if err != nil {
		log.Println(err)
		return ""
	}
	fmt.Println("--------")

	return result
}

func SaveUserInOrderRoom(repo *Implementation, eventId uint, userIdStr string, qUniqueCode string, ctx context.Context) {
	timeNow := time.Now()
	eventIdStr := genEventStr(eventId)

	datas := make(map[string]string)
	datas["time"] = timeNow.String()
	datas["qUniqueCode"] = qUniqueCode
	encData, _ := json.Marshal(datas)
	_, err := repo.redis.HSet(ctx, eventIdStr, userIdStr, encData).Result()
	if err != nil {
		return
	}
}

func (repo *Implementation) GetUserInOrderRoom(userId uint, eventId uint, ctx context.Context) string {
	eventIdStr := genEventStr(eventId)
	userIdStr := strconv.FormatInt(int64(userId), 10)
	result, err := repo.redis.HGet(ctx, eventIdStr, userIdStr).Result()
	if err != nil {
		log.Println(err)
		return ""
	}
	return result
}

func PopUserInOrderRoom(repo *Implementation, userId uint, eventId uint, ctx context.Context) {
	eventIdStr := genEventStr(eventId)
	userIdStr := strconv.FormatInt(int64(userId), 10)
	result, err := repo.redis.HDel(ctx, eventIdStr, userIdStr).Result()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Finish pop user data")
	log.Println(result)
}

func genEventStr(eventId uint) string {
	eventIdStr := "ORmEvent" + strconv.FormatInt(int64(eventId), 10)
	return eventIdStr
}

func genQueueEventStr(eventId uint) string {
	eventIdStr := "QEvent" + strconv.FormatInt(int64(eventId), 10)
	return eventIdStr
}

func CountPeopleInOrderRoom(repo *Implementation, eventId uint, ctx context.Context) int64 {
	eventIdStr := genEventStr(eventId)
	result, err := repo.redis.HLen(ctx, eventIdStr).Result()
	if err != nil {
		log.Println(err)
		return -1
	}
	return result
}
