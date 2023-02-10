package models

import (
	"errors"

	"github.com/google/uuid"
)

var Tokens = make(map[string]string)
var ErrUnauthenticated = errors.New("unauthenticated")

type User struct {
	UUID      uuid.UUID `json:"uuid"`
	FisrtName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}
