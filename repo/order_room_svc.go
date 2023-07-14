package repo

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"ticket-expert/models"
)

var concurrentUser = int64(2)
var helperCountPeople = CountPeopleInOrderRoom
var helperPopQueue = PopWaitingQueue
var helperSave = SaveUserInOrderRoom

func (repo *Implementation) CountTotalPeopleInOrderRoom(eventId uint, ctx context.Context) int64 {
	total := helperCountPeople(repo, eventId, ctx)
	return total
}
// will CountPeopleInOrderRoom, if already slot available, it will get the user in waiting room. 
func (repo *Implementation) CheckOrderRoom(eventId uint, ctx context.Context) []string {
	var qUniqueCodes []string
	for true {
		currUser := helperCountPeople(repo, eventId, ctx)
		log.Printf("Curr capacity : %d", currUser)
		if currUser >= concurrentUser {
			log.Println("Still at max capacity")
			return qUniqueCodes
		}
		nextUserId := helperPopQueue(repo, eventId, ctx)
		if nextUserId == "" {
			log.Println("Queue is empty. Nothing added")
			return qUniqueCodes
		}
		var nextUser models.NewWaitingUser
		json.Unmarshal([]byte(nextUserId), &nextUser)

		log.Printf("moving userId %s to order room", nextUserId)

		qUniqueCodes = append(qUniqueCodes, nextUser.QUniqueCode)
		helperSave(repo, eventId, strconv.Itoa(int(nextUser.UserId)), nextUser.QUniqueCode, ctx)
	}
	return qUniqueCodes
}
