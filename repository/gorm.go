package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormUser struct {
	UUID      uuid.UUID     `json:"uuid" gorm:"primary_key;type:uuid"`
	FisrtName string        `json:"firstName" gorm:"size:50;not null"`
	LastName  string        `json:"lastName" gorm:"size:50;not null"`
	Email     string        `json:"email" gorm:"size:255;not null;unique"`
	Password  string        `json:"password" gorm:"size:250;not null"`
	Accounts  []GormAccount `gorm:"foreignKey:UserUUID"`
}

type GormAccount struct {
	UUID         uuid.UUID `json:"uuid" gorm:"primary_key;type:uuid"`
	IBAN         string    `json:"iban" gorm:"size:250;not null;unique"`
	Balance      uint      `json:"balance" gorm:"not null"`
	UserUUID     uuid.UUID
	Sources      []GormTransaction `gorm:"foreignKey:SourceUUID"`
	Destinations []GormTransaction `gorm:"foreignKey:DestinationUUID"`
}

type GormTransaction struct {
	UUID            uuid.UUID `json:"uuid" gorm:"primary_key;type:uuid"`
	Status          string    `json:"status" gorm:"size:50;not null"`
	SourceUUID      uuid.UUID `gorm:"type:uuid;not null"`
	DestinationUUID uuid.UUID `gorm:"type:uuid;not null"`
	Amount          uint      `gorm:"not null"`
}

func ConnectDataBase() *gorm.DB {

	var DB *gorm.DB
	Dbdriver, ok := os.LookupEnv("DB_DRIVER")
	if !ok {
		log.Fatal("please specify DB_DRIVER")
	}
	DbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		log.Fatal("please specify DB_HOST")
	}
	DbUser, ok := os.LookupEnv("DB_USER")
	if !ok {
		log.Fatal("please specify DB_USER")
	}
	DbPassword, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		log.Fatal("please specify DB_PASSWORD")
	}
	DbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		log.Fatal("please specify DB_NAME")
	}
	DbPort, ok := os.LookupEnv("DB_PORT")
	if !ok {
		log.Fatal("please specify DB_PORT")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DbHost, DbUser, DbPassword, DbName, DbPort)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Println("Cannot connect to database ", Dbdriver)
		log.Fatal("connection error:", err)
	} else {
		log.Println("We are connected to the database ", Dbdriver)
	}

	DB.AutoMigrate(&GormUser{}, &GormAccount{}, &GormTransaction{})
	return DB

}

func ClearData(db *gorm.DB) {
	db.Where("1 = 1").Delete(&GormTransaction{})
	db.Where("1 = 1").Delete(&GormAccount{})
	db.Where("1 = 1").Delete(&GormUser{})
}
