package repository

import (
	"errors"
	"pay/models"

	"github.com/google/uuid"
)

var ErrorCreated = errors.New("user has already created")
var ErrorUnknownUser = errors.New("user does not exist")
var ErrorUnknownAccount = errors.New("account does not exist")
var ErrorUnknownTransaction = errors.New("transaction does not exist")

type TestRepo struct {
	Users       map[uuid.UUID]*models.User
	Accounts    map[uuid.UUID]*models.Account
	Transaction map[uuid.UUID]*models.Transaction
}

func (p *TestRepo) SendTransaction(transactionUUID uuid.UUID) error {
	transaction, ok := p.Transaction[transactionUUID]
	if !ok {
		return ErrorUnknownTransaction
	}
	transaction.Status = StatusSent
	return nil
}

func (t *TestRepo) GetAccountByUUID(uuid uuid.UUID) (*models.Account, error) {
	account, ok := t.Accounts[uuid]
	if !ok {
		return &models.Account{}, ErrorUnknownAccount
	}
	return account, nil
}

func (t *TestRepo) GetTransactionForAccount(accountUUID uuid.UUID) ([]models.Transaction, error) {
	transactions := make([]models.Transaction, 0)
	for _, tr := range t.Transaction {
		if tr.SourceUUID == accountUUID || tr.DestinationUUID == accountUUID {
			transaction := models.Transaction{
				UUID:            tr.UUID,
				Status:          tr.Status,
				SourceUUID:      tr.SourceUUID,
				DestinationUUID: tr.DestinationUUID,
				Amount:          tr.Amount,
			}
			transactions = append(transactions, transaction)
		}
	}
	return transactions, nil
}

func (t *TestRepo) CreateTransaction(transaction models.Transaction) error {
	_, ok := t.Transaction[transaction.UUID]
	if !ok {
		t.Transaction[transaction.UUID] = &transaction
		return nil
	}
	return ErrorCreated
}

func NewTestRepo() TestRepo {
	users := make(map[uuid.UUID]*models.User)
	accounts := make(map[uuid.UUID]*models.Account)
	transaction := make(map[uuid.UUID]*models.Transaction)
	return TestRepo{
		Users:       users,
		Accounts:    accounts,
		Transaction: transaction,
	}
}

func (t *TestRepo) GetAccountsForUser(userUUID uuid.UUID) ([]models.Account, error) {
	accounts := make([]models.Account, 0)
	for _, account := range t.Accounts {
		if account.UserUUID == userUUID {
			accounts = append(accounts, *account)
		}
	}
	return accounts, nil
}

func (p *TestRepo) GetUserByUUID(uuid uuid.UUID) (*models.User, error) {
	user, ok := p.Users[uuid]
	if !ok {
		return &models.User{}, ErrorUnknownUser
	}
	return user, nil
}

func (t *TestRepo) GetUserByEmail(email string) (*models.User, error) {
	for _, user := range t.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return &models.User{}, ErrorUnknownUser
}

func (t *TestRepo) CreateAccount(account *models.Account) error {
	_, ok := t.Accounts[account.UUID]
	if !ok {
		t.Accounts[account.UUID] = account
		user, err := t.GetUserByUUID(account.UserUUID)
		if err != nil {
			return err
		}
		user.Accounts = append(user.Accounts, *account)
		return nil
	}
	return ErrorCreated
}

func (t *TestRepo) CreateUser(user *models.User) error {
	_, ok := t.Users[user.UUID]
	if !ok {
		err := t.CheckIfExistUser(user)
		if err != nil {
			return err
		}
		t.Users[user.UUID] = user
		return nil
	}
	return ErrorCreated
}

func (t *TestRepo) CheckIfExistUser(user *models.User) error {
	for _, us := range t.Users {
		if us.Email == user.Email {
			return ErrorCreated
		}
	}
	return nil
}
