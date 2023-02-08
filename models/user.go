package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/matthewhartstonge/argon2"
)

var Tokens = make(map[string]string)
var ErrUnauthenticated = errors.New("unauthenticated")

func GetEmail(token string) (string, bool) {
	email, ok := Tokens[token]
	return email, ok
}

type User struct {
	gorm.Model
	UUID      uuid.UUID `json:"uuid" gorm:"type:uuid"`
	FisrtName string    `json:"firstName" gorm:"size:50;not null"`
	LastName  string    `json:"lastName" gorm:"size:50;not null"`
	Email     string    `json:"email" gorm:"size:255;not null;unique"`
	Password  string    `json:"password" gorm:"size:250;not null"`
}

func (u *User) CreateUser(DB *gorm.DB) (*User, error) {
	err := DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// Function called automatically before CreateUser()
func (u *User) BeforeCreate() error {
	argon := argon2.DefaultConfig()

	hashedPasword, err := argon.HashEncoded([]byte(u.Password))
	if err != nil {
		return err
	}
	u.UUID = uuid.New()
	u.Password = string(hashedPasword)
	u.FisrtName = strings.TrimSpace(u.FisrtName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(u.Email)
	return nil
}

func LoginCheck(DB *gorm.DB, email string, password string) (string, error) {
	u := User{}
	err := DB.Model(User{}).Where("email = ?", email).Take(&u).Error
	if err != nil {
		return "", ErrUnauthenticated
	}
	ok, err := argon2.VerifyEncoded([]byte(password), []byte(u.Password))
	if err != nil {
		return "", err
	}
	if !ok {
		return "", ErrUnauthenticated
	}
	token, err := randToken(32)
	if err != nil {
		return "", err
	}
	Tokens[token] = email
	return token, nil
}

func randToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
