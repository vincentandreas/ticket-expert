package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/jcuga/golongpoll"
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

var store *sessions.CookieStore

func init() {
	store = sessions.NewCookieStore([]byte("your-secret-key"))
}

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

	lpMngr, err := golongpoll.StartLongpoll(golongpoll.Options{
		LoggingEnabled:                 true,
		DeleteEventAfterFirstRetrieval: true,
	})

	implementObj := repo.NewImplementation(db, redis)
	//if err := gocron.Every(15).Second().Do(implementObj.CheckBookingPeriodically, context.TODO()); err != nil {
	//	panic(err)
	//	return
	//}
	//<-gocron.Start()
	h := controller.NewBaseHandler(implementObj, lpMngr, store)

	log.Fatal(http.ListenAndServe(":10000", controller.HandleRequests(h, lpMngr)))

}
