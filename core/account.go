package core

import (
	"payment/models"

	"github.com/google/uuid"
)

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
