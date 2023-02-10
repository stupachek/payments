package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"

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
	ID        int
	FisrtName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (u *User) Register() error {
	argon := argon2.DefaultConfig()

	hashedPasword, err := argon.HashEncoded([]byte(u.Password))
	u.Password = ""
	if err != nil {
		return err
	}
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
