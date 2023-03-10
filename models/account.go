package models

import (
	"github.com/google/uuid"
)

type Account struct {
	UUID     uuid.UUID `json:"uuid"`
	IBAN     string    `json:"iban"`
	Balance  uint      `json:"balance"`
	UserUUID uuid.UUID `json:"user_uuid"`
	Status   string    `json:"status"`
}
