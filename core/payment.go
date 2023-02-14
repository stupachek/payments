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
var ErrUnknownAccount = errors.New("unknown account")
var ErrInsufficientFunds = errors.New("insufficient funds")

type Transaction struct {
	UserUUID        uuid.UUID
	SourseUUID      uuid.UUID
	DestinationUUID uuid.UUID
	Amount          uint
}

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
	err = p.UserRepo.CreateUser(user)
	return err
}

func (p *PaymentSystem) LoginCheck(email string, password string) (string, error) {
	u, err := p.UserRepo.GetUserByEmail(email)
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

func (p *PaymentSystem) NewAccount(userUUID uuid.UUID) (models.Account, error) {
	user, err := p.UserRepo.GetUserByUUID(userUUID)
	if err != nil {
		return models.Account{}, err
	}
	account := models.Account{}
	account.UserUUID = user.UUID
	account.IBAN, err = randToken(29)
	if err != nil {
		return models.Account{}, err
	}
	account.UUID, err = uuid.NewRandom()
	if err != nil {
		return models.Account{}, err
	}
	err = p.UserRepo.CreateAccount(&account)
	if err != nil {
		return models.Account{}, err
	}
	return account, err
}

func (p *PaymentSystem) NewTransaction(tr Transaction) (models.Transaction, error) {
	user, err := p.UserRepo.GetUserByUUID(tr.UserUUID)
	if err != nil {
		return models.Transaction{}, err
	}
	source, err := checkAccountExists(user.Accounts, tr.SourseUUID)
	if err != nil {
		return models.Transaction{}, err
	}
	if err := p.checkAmount(source, tr.Amount); err != nil {
		return models.Transaction{}, err
	}
	transaction := models.Transaction{
		Status:          "prepared",
		SourceUUID:      tr.SourseUUID,
		DestinationUUID: tr.DestinationUUID,
		Amount:          tr.Amount,
	}
	transaction.UUID, err = uuid.NewRandom()
	if err != nil {
		return models.Transaction{}, err
	}
	err = p.UserRepo.CreateTransaction(transaction)
	if err != nil {
		return models.Transaction{}, err
	}
	return transaction, nil
}

func checkAccountExists(accounts []models.Account, accountUUID uuid.UUID) (models.Account, error) {
	for _, acc := range accounts {
		if acc.UUID == accountUUID {
			return acc, nil
		}
	}
	return models.Account{}, ErrUnknownAccount
}

func (p *PaymentSystem) checkAmount(account models.Account, amount uint) error {
	if amount <= account.Balance {
		return nil
	}
	return ErrInsufficientFunds
}

func (p *PaymentSystem) GetAccounts(userUUID uuid.UUID) ([]models.Account, error) {
	return p.UserRepo.GetAccounts(userUUID)
}
