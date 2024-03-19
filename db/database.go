package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

var DB *gorm.DB

func InitDB() {
	var dbHost, dbUser, dbPassword, dbName, dbPort string

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		dbHost = "localhost"
		dbUser = "postgres"
		dbPassword = "POSTGRES_PASSWORD"
		dbName = "ginstagram"
		dbPort = "5432"
	} else {
		dbHost = os.Getenv("DB_HOST")
		dbUser = os.Getenv("DB_USER")
		dbPassword = os.Getenv("DB_PASSWORD")
		dbName = os.Getenv("DB_NAME")
		dbPort = os.Getenv("DB_PORT")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Singapore",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db
}
