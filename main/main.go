// swagger:meta
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
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Event{})
	db.AutoMigrate(&models.EventDetail{})
	db.AutoMigrate(&models.PurchasedTicket{})
	db.AutoMigrate(&models.BookingTicket{})
	db.AutoMigrate(&models.BookingDetail{})
}

func dbSetup() (*gorm.DB, *redis.Client, error) {
	username := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DBNAME")
	host := os.Getenv("POSTGRES_HOST")
	dbParams := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s", username, password, dbName, host)
	conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbParams,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWD"),
		DB:       0,
	})
	return conn, redisClient, err
}

// @title           Ticket Expert API
// @version         2.0
// @description     Used for ordering ticket.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
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
	//if err := gocron.Every(15).Second().Do(implementObj.CheckBookingPeriod, context.TODO()); err != nil {
	//	panic(err)
	//	return
	//}
	//<-gocron.Start()
	h := controller.NewBaseHandler(implementObj, lpMngr, store)

	log.Fatal(http.ListenAndServe(":10000", controller.HandleRequests(h, lpMngr)))

}
