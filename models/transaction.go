package models

import (
	"github.com/google/uuid"
)

type Transaction struct {
	ID            uint
	UUID          uuid.UUID `json:"uuid"`
	Status        string    `json:"status"`
	SourceId      uint      `json:"source_id"`
	DestinationId uint      `json:"destination_id"`
	Amount        uint      `json:"money"`
}
