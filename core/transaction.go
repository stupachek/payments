package core

import (
	"errors"
	"payment/models"
	"payment/repository"

	"github.com/google/uuid"
)

const (
	SENT = "sent"
)

var (
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
	transactionModel, err := p.Repo.GetTransactionByUUID(transaction.UUID)
	if err != nil {
		return models.Transaction{}, err
	}
	return *transactionModel, nil
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
			err = repo.UpdateStatusTransaction(transactionUUID, SENT)
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
