package core

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"payment/models"
	"strings"

	"github.com/google/uuid"
	"github.com/matthewhartstonge/argon2"
	"gorm.io/gorm"
)

const (
	USER        = "user"
	ADMIN       = "admin"
	EMAIN_ADMIN = "admin@admin.admin"
)

var (
	Tokens              = make(map[string]string)
	ErrUnauthenticated  = errors.New("unauthenticated")
	ErrPermissionDenied = errors.New("permission denied")
	ErrUserBlocked      = errors.New("user is blocked")
	ErrUserActive       = errors.New("user is active")
	ErrBadRequest       = errors.New("bad request")
)

type LoginReturn struct {
	UUID  uuid.UUID
	Token string
}

func newPassword(password string) (string, error) {
	argon := argon2.DefaultConfig()

	hashedPasword, err := argon.HashEncoded([]byte(password))
	password = ""
	if err != nil {
		return "", err
	}
	return string(hashedPasword), nil

}

func (p *PaymentSystem) Register(user *models.User) error {
	var err error
	user.Password, err = newPassword(user.Password)
	if err != nil {
		return err
	}
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
	if u.Status == BLOCKED {
		return LoginReturn{}, ErrUserBlocked
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
	admin, err := p.Repo.GetUserByEmail(EMAIN_ADMIN)
	password := os.Getenv("PAYMENT_ADMIN_PASSWORD")
	if strings.Contains(err.Error(), "duplicate key value") {
		password, err = newPassword(password)
		if err != nil {
			return err
		}
		p.Repo.UpdatePassword(admin.UUID, password)
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	user := &models.User{
		FisrtName: "admin",
		LastName:  "admin",
		Email:     EMAIN_ADMIN,
		Password:  password,
		Role:      ADMIN,
		Status:    ACTIVE,
	}
	err = p.Register(user)
	return err
}

func (p *PaymentSystem) BlockUser(userUUID uuid.UUID) error {
	ok, err := p.IsBlockedUser(userUUID)
	if err != nil {
		return ErrBadRequest
	}
	if !ok {
		return ErrUserBlocked
	}
	err = p.Repo.UpdateStatusUser(userUUID, BLOCKED)
	if err != nil {
		return ErrBadRequest
	}
	return nil
}
func (p *PaymentSystem) UnblockUser(userUUID uuid.UUID) error {
	ok, err := p.IsBlockedUser(userUUID)
	if err != nil {
		return ErrBadRequest
	}
	if ok {
		return ErrUserActive
	}
	err = p.Repo.UpdateStatusUser(userUUID, ACTIVE)
	if err != nil {
		return ErrBadRequest
	}
	return nil
}

func (p *PaymentSystem) IsBlockedUser(userUUID uuid.UUID) (bool, error) {
	user, err := p.Repo.GetUserByUUID(userUUID)
	if err != nil {
		return false, err
	}
	return user.Status == BLOCKED, nil
}
