package repo

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"ticket-expert/models"
)

var concurrentUser = int64(2)

func (repo *Implementation) CheckOrderRoom(eventId uint, ctx context.Context) []string {
	currUser := repo.CountPeopleInOrderRoom(eventId, ctx)
	log.Printf("Curr capacity : %d", currUser)
	var qUniqueCodes []string
	for true {

		if currUser >= concurrentUser {
			log.Println("Still at max capacity")
			return qUniqueCodes
		}
		nextUserId := repo.PopWaitingQueue(eventId, ctx)
		if nextUserId == "" {
			log.Println("Queue is empty. Nothing added")
			return qUniqueCodes
		}
		var nextUser models.NewWaitingUser
		json.Unmarshal([]byte(nextUserId), &nextUser)

		log.Printf("moving userId %s to order room", nextUserId)

		qUniqueCodes = append(qUniqueCodes, nextUser.QUniqueCode)
		repo.SaveUserInOrderRoom(eventId, strconv.Itoa(int(nextUser.UserId)), nextUser.QUniqueCode, ctx)
	}
	return qUniqueCodes
}
