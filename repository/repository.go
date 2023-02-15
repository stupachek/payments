package repository

import (
	"errors"
	"pay/models"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var ErrorCreated = errors.New("user has already created")
var ErrorUnknownUser = errors.New("user does not exist")
var ErrorUnknownAccount = errors.New("account does not exist")

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUUID(uuid uuid.UUID) (*models.User, error)
	CreateAccount(account *models.Account) error
	CreateTransaction(transaction models.Transaction) error
	GetAccounts(userUUID uuid.UUID) ([]models.Account, error)
	GetTransactionForAccount(accountUUID uuid.UUID) ([]models.Transaction, error)
}

type PostgresRepo struct {
	DB *gorm.DB
}

func (p *PostgresRepo) CreateTransaction(transaction models.Transaction) error {
	gormTransaction := GormTransaction{
		UUID:            transaction.UUID,
		Status:          transaction.Status,
		SourceUUID:      transaction.SourceUUID,
		DestinationUUID: transaction.DestinationUUID,
		Amount:          transaction.Amount,
	}
	err := p.DB.Create(&gormTransaction).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepo) GetTransactionForAccount(accountUUID uuid.UUID) ([]models.Transaction, error) {
	var gormTransaction []GormTransaction
	result := p.DB.Model(GormTransaction{}).Find(&gormTransaction).Where("SourceUUID = ? OR DestinationUUID = ?", accountUUID, accountUUID)
	if result.Error != nil {
		return []models.Transaction{}, result.Error
	}
	modelTransaction := p.fromGormToModelTransaction(gormTransaction)
	return modelTransaction, nil
}

func (p *PostgresRepo) fromGormToModelAccount(accounts []GormAccount) []models.Account {
	modelAccounts := make([]models.Account, len(accounts))
	for i, acc := range accounts {
		modelAccounts[i] = models.Account{
			UUID:     acc.UUID,
			IBAN:     acc.IBAN,
			Balance:  acc.Balance,
			UserUUID: acc.UserUUID,
		}
	}

	return modelAccounts
}

func (p *PostgresRepo) fromGormToModelTransaction(transactions []GormTransaction) []models.Transaction {
	modelTransaction := make([]models.Transaction, len(transactions))
	for i, tr := range transactions {
		source, err := p.GetAccountByUUID(tr.SourceUUID)
		if err != nil {
			return nil
		}
		destination, err := p.GetAccountByUUID(tr.DestinationUUID)
		if err != nil {
			return nil
		}
		modelTransaction[i] = models.Transaction{
			UUID:            tr.UUID,
			Status:          tr.Status,
			SourceUUID:      source.UUID,
			DestinationUUID: destination.UUID,
			Amount:          tr.Amount,
		}
	}
	return modelTransaction
}

func (p *PostgresRepo) GetAccounts(userUUID uuid.UUID) ([]models.Account, error) {
	var gormAccounts []GormAccount
	result := p.DB.Model(GormAccount{}).Find(&gormAccounts).Where("UserUUID = ?", userUUID).Preload("Sources").Preload("Destinations")
	if result.Error != nil {
		return []models.Account{}, result.Error
	}
	modelAccounts := p.fromGormToModelAccount(gormAccounts)
	return modelAccounts, nil

}

func (p *PostgresRepo) CreateAccount(account *models.Account) error {
	gormAcc := GormAccount{
		UUID:     account.UUID,
		IBAN:     account.IBAN,
		Balance:  account.Balance,
		UserUUID: account.UserUUID,
	}
	err := p.DB.Create(&gormAcc).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepo) GetAccountByUUID(uuid uuid.UUID) (*models.Account, error) {
	gormAccount := GormAccount{}
	err := p.DB.Model(GormAccount{}).Where("UUID = ?", uuid).Take(&gormAccount)
	if err != nil {
		return &models.Account{}, nil
	}
	account := models.Account{
		UUID:     gormAccount.UUID,
		IBAN:     gormAccount.IBAN,
		Balance:  gormAccount.Balance,
		UserUUID: gormAccount.UserUUID,
	}
	return &account, nil
}

func (p *PostgresRepo) GetUserByUUID(uuid uuid.UUID) (*models.User, error) {
	userGorm := GormUser{}
	err := p.DB.Model(GormUser{}).Where("UUID = ?", uuid).Preload("Accounts").Take(&userGorm).Error
	if err != nil {
		return &models.User{}, err
	}
	user := models.User{
		UUID:      userGorm.UUID,
		FisrtName: userGorm.FisrtName,
		LastName:  userGorm.LastName,
		Email:     userGorm.Email,
		Password:  userGorm.Password,
		Accounts:  p.fromGormToModelAccount(userGorm.Accounts),
	}
	return &user, nil
}

func (p *PostgresRepo) GetUserByEmail(email string) (*models.User, error) {
	userGorm := GormUser{}
	err := p.DB.Model(GormUser{}).Where("email = ?", email).Preload("Accounts").Take(&userGorm).Error
	if err != nil {
		return &models.User{}, err
	}
	user := models.User{
		UUID:      userGorm.UUID,
		FisrtName: userGorm.FisrtName,
		LastName:  userGorm.LastName,
		Email:     userGorm.Email,
		Password:  userGorm.Password,
		Accounts:  p.fromGormToModelAccount(userGorm.Accounts),
	}
	return &user, nil
}

func (p *PostgresRepo) CreateUser(user *models.User) error {

	gormU := GormUser{
		UUID:      user.UUID,
		FisrtName: user.FisrtName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	}
	err := p.DB.Create(&gormU).Error
	if err != nil {
		return err
	}
	return nil
}

func NewGormUserRepo(DB *gorm.DB) UserRepository {
	return &PostgresRepo{
		DB: DB,
	}
}
