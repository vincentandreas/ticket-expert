package repo

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"strconv"
	"ticket-expert/models"
)

var concurrentUser = int64(2)

func (repo *Implementation) CheckOrderRoom(eventId uint, ctx context.Context) {
	currUser := repo.CountPeopleInOrderRoom(eventId, ctx)
	log.Printf("Curr capacity : %d", currUser)
	for true {
		uniqueId := uuid.New().String()
		if currUser >= concurrentUser {
			log.Println("Still at max capacity")
			return
		}
		nextUserId := repo.PopWaitingQueue(eventId, ctx)
		if nextUserId == "" {
			log.Println("Queue is empty. Nothing added")
			return
		}
		var nextUser models.NewWaitingUser
		json.Unmarshal([]byte(nextUserId), &nextUser)

		log.Printf("moving userId %s to order room", nextUserId)
		repo.SaveUserInOrderRoom(eventId, strconv.Itoa(int(nextUser.UserId)), uniqueId, ctx)
	}
}
