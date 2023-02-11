package core

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"pay/models"
	"pay/repository"
	"strings"

	"github.com/google/uuid"
	"github.com/matthewhartstonge/argon2"
)

var Tokens = make(map[string]string)
var ErrUnauthenticated = errors.New("unauthenticated")

func GetEmail(token string) (string, bool) {
	email, ok := Tokens[token]
	return email, ok
}

type PaymentSystem struct {
	UserRepo repository.UserRepository
}

func NewPaymentSystem(userRepo repository.UserRepository) PaymentSystem {
	return PaymentSystem{
		UserRepo: userRepo,
	}
}

func (p PaymentSystem) Register(user *models.User) error {
	argon := argon2.DefaultConfig()

	hashedPasword, err := argon.HashEncoded([]byte(user.Password))
	user.Password = ""
	if err != nil {
		return err
	}
	user.Password = string(hashedPasword)
	user.UUID, err = uuid.NewRandom()
	if err != nil {
		return err
	}
	user.FisrtName = strings.TrimSpace(user.FisrtName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Email = strings.TrimSpace(user.Email)
	err = p.UserRepo.CreateUser(user)
	return err
}

func (p PaymentSystem) LoginCheck(email string, password string) (string, error) {
	u, err := p.UserRepo.GetUserByEmail(email)
	if err != nil {
		return "", err
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

func (p PaymentSystem) CheckToken(UUID uuid.UUID, token string) error {
	user, err := p.UserRepo.GetUserByUUID(UUID)
	if err != nil {
		return ErrUnauthenticated
	}
	if token == "" {
		return ErrUnauthenticated
	}
	email, ok := GetEmail(token)
	if !ok {
		return ErrUnauthenticated
	}
	if email != user.Email {
		return ErrUnauthenticated
	}
	return nil
}
