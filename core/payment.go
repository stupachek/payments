package core

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"payment/models"
	"payment/repository"
	"strings"

	"github.com/google/uuid"
	"github.com/matthewhartstonge/argon2"
)

var (
	StatusSent           = "sent"
	Tokens               = make(map[string]string)
	ErrUnauthenticated   = errors.New("unauthenticated")
	ErrUnknownAccount    = errors.New("unknown account")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrWrongDestination  = errors.New("source equals destination")
)

type Transaction struct {
	UserUUID        uuid.UUID
	SourceUUID      uuid.UUID
	DestinationUUID uuid.UUID
	Amount          uint
}

func GetEmail(token string) (string, bool) {
	email, ok := Tokens[token]
	return email, ok
}

type PaymentSystem struct {
	Repo repository.Repository
}

func NewPaymentSystem(userRepo repository.Repository) PaymentSystem {
	return PaymentSystem{
		Repo: userRepo,
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

func (p *PaymentSystem) NewAccount(userUUID uuid.UUID) (models.Account, error) {
	user, err := p.Repo.GetUserByUUID(userUUID)
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
	err = p.Repo.CreateAccount(&account)
	if err != nil {
		return models.Account{}, err
	}
	return account, err
}

func (p *PaymentSystem) NewTransaction(tr Transaction) (models.Transaction, error) {
	if tr.SourceUUID == tr.DestinationUUID {
		return models.Transaction{}, ErrWrongDestination
	}
	err := p.checkAmount(tr.SourceUUID, tr.Amount)
	if err != nil {
		return models.Transaction{}, err
	}
	transaction := models.Transaction{
		Status:          "prepared",
		SourceUUID:      tr.SourceUUID,
		DestinationUUID: tr.DestinationUUID,
		Amount:          tr.Amount,
	}
	transaction.UUID, err = uuid.NewRandom()
	if err != nil {
		return models.Transaction{}, err
	}
	err = p.Repo.CreateTransaction(transaction)
	if err != nil {
		return models.Transaction{}, err
	}
	return transaction, nil
}

func (p *PaymentSystem) CheckAccountExists(userUUID, accountUUID uuid.UUID) error {
	account, err := p.Repo.GetAccountByUUID(accountUUID)
	if err != nil {
		return err
	}
	if account.UserUUID != userUUID {
		return ErrUnknownAccount
	}
	return nil
}

func (p *PaymentSystem) checkAmount(accountUUID uuid.UUID, amount uint) error {
	account, err := p.Repo.GetAccountByUUID(accountUUID)
	if err != nil {
		return err
	}
	if amount <= account.Balance {
		return nil
	}
	return ErrInsufficientFunds
}

func (p *PaymentSystem) GetAccounts(userUUID uuid.UUID, query models.QueryParams) ([]models.Account, error) {
	return p.Repo.GetAccountsForUser(userUUID, query)
}

func (p *PaymentSystem) GetTransactions(accountUUID uuid.UUID, query models.QueryParams) ([]models.Transaction, error) {
	return p.Repo.GetTransactionForAccount(accountUUID, query)
}

func (p *PaymentSystem) SendTransaction(transactionUUID uuid.UUID) (models.Transaction, error) {
	err := p.Repo.Transaction(
		func(repo repository.Repository) error {
			transaction, err := repo.GetTransactionByUUID(transactionUUID)
			if err != nil {
				return err
			}
			err = p.checkAmount(transaction.SourceUUID, transaction.Amount)
			if err != nil {
				return err
			}
			err = repo.DecBalance(transaction.SourceUUID, transaction.Amount)
			if err != nil {
				return err
			}
			err = repo.IncBalance(transaction.DestinationUUID, transaction.Amount)
			if err != nil {
				return err
			}
			err = repo.UpdateStatus(transactionUUID, StatusSent)
			return err
		})
	if err != nil {
		return models.Transaction{}, err
	}
	tr, err := p.Repo.GetTransactionByUUID(transactionUUID)
	if err != nil {
		return models.Transaction{}, err
	}
	return *tr, nil
}

func (p *PaymentSystem) AddMoney(accountUUID uuid.UUID, amount uint) (models.Account, error) {
	p.Repo.IncBalance(accountUUID, amount)
	account, err := p.Repo.GetAccountByUUID(accountUUID)
	if err != nil {
		return models.Account{}, err
	}
	return *account, nil
}

func (p *PaymentSystem) ShowBalance(accountUUID uuid.UUID) (uint, error) {
	account, err := p.Repo.GetAccountByUUID(accountUUID)
	if err != nil {
		return 0, err
	}
	return account.Balance, nil

}

func (p *PaymentSystem) GetAccount(accountUUID uuid.UUID) (models.Account, error) {
	account, err := p.Repo.GetAccountByUUID(accountUUID)
	return *account, err
}
