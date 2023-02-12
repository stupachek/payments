package models

import (
	"github.com/google/uuid"
)

type Transaction struct {
	UUID          uuid.UUID `json:"uuid"`
	Status        string    `json:"status"`
	SourceId      uint      `json:"source_id"`
	DestinationId uint      `json:"destination_id"`
}
