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
	result, err := repo.redis.LPush(ctx, key, wuserJson).Result()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(result)
}

func (repo *Implementation) PopWaitingQueue(eventId uint, ctx context.Context) string {
	key := genQueueEventStr(eventId)
	result, err := repo.redis.RPop(ctx, key).Result()
	if err != nil {
		log.Println(err)
		return ""
	}
	fmt.Println("--------")
	fmt.Println(result)
	return result
}

func (repo *Implementation) SaveUserInOrderRoom(eventId uint, userIdStr string, qUniqueCode string, ctx context.Context) {
	timeNow := time.Now()
	eventIdStr := genEventStr(eventId)

	datas := make(map[string]string)
	datas["time"] = timeNow.String()
	datas["queueUniqueCode"] = qUniqueCode
	encData, _ := json.Marshal(datas)
	result, err := repo.redis.HSet(ctx, eventIdStr, userIdStr, encData).Result()
	if err != nil {
		return
	}
	fmt.Println(result)
	fmt.Println("->>>>")
}

func (repo *Implementation) GetUserInOrderRoom(userId uint, eventId uint, ctx context.Context) string {
	eventIdStr := genEventStr(eventId)
	userIdStr := strconv.FormatInt(int64(userId), 10)
	result, err := repo.redis.HGet(ctx, eventIdStr, userIdStr).Result()
	if err != nil {
		log.Println(err)
		return ""
	}
	fmt.Println(result)
	fmt.Println("->>>>")
	return result
}

func (repo *Implementation) PopUserInOrderRoom(userId uint, eventId uint, ctx context.Context) {
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

func (repo *Implementation) CountPeopleInOrderRoom(eventId uint, ctx context.Context) int64 {
	eventIdStr := genEventStr(eventId)
	result, err := repo.redis.HLen(ctx, eventIdStr).Result()
	if err != nil {
		log.Println(err)
		return -1
	}
	return result
}
