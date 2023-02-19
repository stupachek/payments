package models

import (
	"github.com/google/uuid"
)

type Account struct {
	UUID     uuid.UUID `json:"uuid"`
	IBAN     string    `json:"iban"`
	Balance  uint      `json:"balance"`
	UserUUID uuid.UUID `json:"user_uuid"`
}

type QueryParams struct {
	Limit  uint `json:"limit"`
	Offset uint `json:"offset"`
}
