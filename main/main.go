package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

func dbSetup() (*gorm.DB, error) {
	username := os.Getenv("postgres_username")
	password := os.Getenv("postgres_password")
	dbName := os.Getenv("postgres_dbname")
	host := os.Getenv("postgres_host")
	dbParams := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s", username, password, dbName, host)
	conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbParams,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	return conn, err
}
func main() {
	db, err := dbSetup()
	if err != nil {
		panic(err)
	}
	doMigration(db)

	implementObj := repo.NewImplementation(db)

	h := controller.NewBaseHandler(implementObj)

	log.Fatal(http.ListenAndServe(":10000", controller.HandleRequests(h)))

}
