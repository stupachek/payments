package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	UUID            uuid.UUID `json:"uuid"`
	Status          string    `json:"status"`
	SourceUUID      uuid.UUID `json:"source_uuid"`
	DestinationUUID uuid.UUID `json:"destination_uuid"`
	Amount          uint      `json:"amount"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
