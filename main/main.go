package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"ticket-expert/controller"
	"ticket-expert/models"
	"ticket-expert/repo"
)

func doMigration(db *gorm.DB) {
	db.AutoMigrate(&models.Promotor{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Event{})
	db.AutoMigrate(&models.EventDetail{})
	db.AutoMigrate(&models.PurchasedTicket{})
	db.AutoMigrate(&models.BookingTicket{})
	db.AutoMigrate(&models.BookingDetail{})
}

func dbSetup() (*gorm.DB, *redis.Client, error) {
	username := os.Getenv("postgres_username")
	password := os.Getenv("postgres_password")
	dbName := os.Getenv("postgres_dbname")
	host := os.Getenv("postgres_host")
	dbParams := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s", username, password, dbName, host)
	conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbParams,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis_host"),
		Password: os.Getenv("redis_passwd"),
		DB:       0,
	})
	return conn, redisClient, err
}
func main() {
	db, redis, err := dbSetup()
	if err != nil {
		panic(err)
	}
	doMigration(db)

	implementObj := repo.NewImplementation(db, redis)
	//if err := gocron.Every(15).Second().Do(implementObj.CheckBookingPeriod, context.TODO()); err != nil {
	//	panic(err)
	//	return
	//}
	//<-gocron.Start()
	h := controller.NewBaseHandler(implementObj)
	log.Fatal(http.ListenAndServe(":10000", controller.HandleRequests(h)))

}
