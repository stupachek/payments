package core

import (
	"crypto/rand"
	"encoding/hex"
	"payment/models"
	"strings"

	"github.com/google/uuid"
	"github.com/matthewhartstonge/argon2"
)

func (p *PaymentSystem) Register(user *models.User) error {
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
	user.Role = USER
	err = p.Repo.CreateUser(user)
	return err
}

func (p *PaymentSystem) LoginCheck(email string, password string) (string, error) {
	u, err := p.Repo.GetUserByEmail(email)
	if err != nil {
		return "", ErrUnauthenticated
	}
	ok, err := argon2.VerifyEncoded([]byte(password), []byte(u.Password))
	if err != nil {
		return "", ErrUnauthenticated
	}
	if !ok {
		return "", ErrUnauthenticated
	}
	token, err := randToken(32)
	if err != nil {
		return "", ErrUnauthenticated
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
func (p *PaymentSystem) CheckToken(UUID uuid.UUID, token string) error {
	user, err := p.Repo.GetUserByUUID(UUID)
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
