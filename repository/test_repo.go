package repository

import (
	"errors"
	"payment/models"

	"github.com/google/uuid"
)

var ErrorCreated = errors.New("user has already created")
var ErrorUnknownUser = errors.New("user does not exist")
var ErrorUnknownAccount = errors.New("account does not exist")
var ErrorUnknownTransaction = errors.New("transaction does not exist")

type TestRepo struct {
	Users        map[uuid.UUID]*models.User
	Accounts     map[uuid.UUID]*models.Account
	Transactions map[uuid.UUID]*models.Transaction
}

func (t *TestRepo) Transaction(callback func(repo Repository) error) error {
	return callback(t)
}

func (t *TestRepo) UpdateStatus(transactionUUID uuid.UUID, status string) error {
	transaction, ok := t.Transactions[transactionUUID]
	if !ok {
		return ErrorUnknownAccount
	}
	transaction.Status = status
	return nil
}

func (t *TestRepo) DecBalance(accountUUID uuid.UUID, amount uint) error {
	account, ok := t.Accounts[accountUUID]
	if !ok {
		return ErrorUnknownAccount
	}
	account.Balance = account.Balance - amount
	return nil
}

func (t *TestRepo) IncBalance(accountUUID uuid.UUID, amount uint) error {
	account, ok := t.Accounts[accountUUID]
	if !ok {
		return ErrorUnknownAccount
	}
	account.Balance = account.Balance + amount
	return nil
}

func (t *TestRepo) GetAccountByUUID(uuid uuid.UUID) (*models.Account, error) {
	account, ok := t.Accounts[uuid]
	if !ok {
		return &models.Account{}, ErrorUnknownAccount
	}
	return account, nil
}

func (t *TestRepo) GetTransactionForAccount(accountUUID uuid.UUID, query models.QueryParams) ([]models.Transaction, error) {
	transactions := make([]models.Transaction, 0)
	for _, tr := range t.Transactions {
		if tr.SourceUUID == accountUUID || tr.DestinationUUID == accountUUID {
			transactions = append(transactions, *tr)
		}
	}
	return transactions, nil
}
func (t *TestRepo) GetTransactionByUUID(transactionUUID uuid.UUID) (*models.Transaction, error) {
	transaction, ok := t.Transactions[transactionUUID]
	if !ok {
		return &models.Transaction{}, ErrorUnknownTransaction
	}
	return transaction, nil
}
func (t *TestRepo) CreateTransaction(transaction models.Transaction) error {
	_, ok := t.Transactions[transaction.UUID]
	if !ok {
		t.Transactions[transaction.UUID] = &transaction
		return nil
	}
	return ErrorCreated
}

func NewTestRepo() TestRepo {
	users := make(map[uuid.UUID]*models.User)
	accounts := make(map[uuid.UUID]*models.Account)
	transaction := make(map[uuid.UUID]*models.Transaction)
	return TestRepo{
		Users:        users,
		Accounts:     accounts,
		Transactions: transaction,
	}
}

func (t *TestRepo) GetAccountsForUser(userUUID uuid.UUID, paganition models.QueryParams) ([]models.Account, error) {
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
