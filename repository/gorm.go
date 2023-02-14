package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
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
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	Dbdriver := os.Getenv("DB_DRIVER")
	DbHost := os.Getenv("DB_HOST")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DbHost, DbUser, DbPassword, DbName, DbPort)
	DB, err = gorm.Open(Dbdriver, dsn)

	if err != nil {
		log.Println("Cannot connect to database ", Dbdriver)
		log.Fatal("connection error:", err)
	} else {
		log.Println("We are connected to the database ", Dbdriver)
	}

	DB.AutoMigrate(&GormUser{}, &GormAccount{}, &GormTransaction{})
	return DB

}
