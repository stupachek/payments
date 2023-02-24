package core

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"payment/models"
	"strings"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/matthewhartstonge/argon2"
)

const (
	USER  = "user"
	ADMIN = "admin"
)

var (
	Tokens              = make(map[string]string)
	ErrUnauthenticated  = errors.New("unauthenticated")
	ErrPermissionDenied = errors.New("permission denied")
)

type LoginReturn struct {
	UUID  uuid.UUID
	Token string
}

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
	err = p.Repo.CreateUser(user)
	return err
}

func (p *PaymentSystem) LoginCheck(email string, password string) (LoginReturn, error) {
	u, err := p.Repo.GetUserByEmail(email)
	if err != nil {
		return LoginReturn{}, ErrUnauthenticated
	}
	ok, err := argon2.VerifyEncoded([]byte(password), []byte(u.Password))
	if err != nil {
		return LoginReturn{}, ErrUnauthenticated
	}
	if !ok {
		return LoginReturn{}, ErrUnauthenticated
	}
	token, err := randToken(32)
	if err != nil {
		return LoginReturn{}, ErrUnauthenticated
	}
	Tokens[token] = email
	loginReturn := LoginReturn{
		UUID:  u.UUID,
		Token: token,
	}
	return loginReturn, nil
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

func (p *PaymentSystem) ChangeRole(adminUUID, userUUID uuid.UUID, role string) error {
	err := p.Repo.UpdateRole(userUUID, role)
	return err
}

func (p *PaymentSystem) CheckAdmin(UUID uuid.UUID) error {
	admin, err := p.Repo.GetUserByUUID(UUID)
	if err != nil {
		return ErrPermissionDenied
	}
	if admin.Role != ADMIN {
		return ErrPermissionDenied
	}
	return nil
}

func (p *PaymentSystem) SetupAdmin() error {
	admin, _ := p.Repo.GetUserByEmail("admin@admin.admin")
	if admin.Email != "" {
		return nil
	}
	godotenv.Load(".env")
	password := os.Getenv("PAYMENT_ADMIN_PASSWORD")
	user := &models.User{
		FisrtName: "admin",
		LastName:  "admin",
		Email:     "admin@admin.admin",
		Password:  password,
		Role:      ADMIN,
	}
	err := p.Register(user)
	return err
}
