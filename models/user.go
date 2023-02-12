package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID        uint      `json:"-"`
	UUID      uuid.UUID `json:"uuid"`
	FisrtName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Accounts  []Account `json:"accounts"`
}
