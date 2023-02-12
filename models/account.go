package models

import (
	"github.com/google/uuid"
)

type Account struct {
	ID           uint          `json:"-"`
	UUID         uuid.UUID     `json:"uuid"`
	IBAN         string        `json:"iban"`
	Balance      float64       `json:"balance"`
	UserId       uint          `json:"user_id"`
	Sources      []Transaction `json:"sources"`
	Destinations []Transaction `json:"destination"`
}
