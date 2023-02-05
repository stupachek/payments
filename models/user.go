package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"html"
	"strings"

	"github.com/google/uuid"
	"github.com/matthewhartstonge/argon2"
	"gorm.io/gorm"
)

var Tokens = make(map[string]string)

func GetToken(email string) (string, bool) {
	tok, ok := Tokens[email]
	return tok, ok
}

type User struct {
	gorm.Model
	// ID        uint   `json:"id" gorm:"primary_key"`
	UUID      uuid.UUID `json:"uuid" gorm:"type:uuid"`
	FisrtName string    `json:"firstName" gorm:"size:50;not null"`
	LastName  string    `json:"lastName" gorm:"size:50;not null"`
	Email     string    `json:"email" gorm:"size:255;not null;unique"`
	Password  string    `json:"password" gorm:"size:250;not null"`
}

func (u *User) CreateUser() (*User, error) {
	err := DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeCreate() error {
	argon := argon2.DefaultConfig()

	hashedPasword, err := argon.HashEncoded([]byte(u.Password))
	if err != nil {
		return err
	}
	u.UUID = uuid.New()
	u.Password = string(hashedPasword)
	u.FisrtName = html.EscapeString(strings.TrimSpace(u.FisrtName))
	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))
	return nil
}

func LoginCheck(email string, password string) (string, error) {
	u := User{}
	err := DB.Model(User{}).Where("email = ?", email).Take(&u).Error
	if err != nil {
		return "", errors.New("unauthenticated")
	}
	ok, err := argon2.VerifyEncoded([]byte(password), []byte(u.Password))
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("unauthenticated")
	}
	token, err := randToken(32)
	if err != nil {
		return "", err
	}
	Tokens[email] = token
	return token, nil
}

func randToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
