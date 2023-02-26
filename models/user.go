package models

import (
	"github.com/google/uuid"
)

type User struct {
	UUID      uuid.UUID `json:"uuid"`
	FisrtName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	Accounts  []Account `json:"accounts"`
}
