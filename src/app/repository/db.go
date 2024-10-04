package repository

import (
	"os"
	"log"
	"gorm.io/gorm"
    "gorm.io/driver/postgres"
)

func Engine() *gorm.DB {
	dsn := "user="+os.Getenv("dbuser")+
		" password="+os.Getenv("dbpassword")+
		" database="+os.Getenv("dbuser")+
		" host="+os.Getenv("dbuser")+
		" port="+os.Getenv("dbport")+
		" sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
        log.Fatal(err)
    }
	return db
}