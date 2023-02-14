package models

import (
	"github.com/google/uuid"
)

type Transaction struct {
	UUID            uuid.UUID `json:"uuid"`
	Status          string    `json:"status"`
	SourceUUID      uuid.UUID `json:"source_uuid"`
	DestinationUUID uuid.UUID `json:"destination_uuid"`
	Amount          uint      `json:"amount"`
}
