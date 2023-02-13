package models

import (
	"github.com/google/uuid"
)

type Account struct {
	ID           uint          `json:"-"`
	UUID         uuid.UUID     `json:"uuid"`
	IBAN         string        `json:"iban"`
	Balance      uint          `json:"balance"`
	UserId       uint          `json:"-"`
	Sources      []Transaction `json:"sources"`
	Destinations []Transaction `json:"destination"`
}
