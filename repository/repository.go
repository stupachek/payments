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
}

type PostgresRepo struct {
	DB *gorm.DB
}

type TestRepo struct {
	Users       map[uuid.UUID]*models.User
	Accounts    map[uuid.UUID]*models.Account
	Transaction map[uuid.UUID]*models.Transaction
}

func (p *PostgresRepo) CreateTransaction(transaction models.Transaction) error {
	source, err := p.GetAccountByUUID(transaction.SourceUUID)
	if err != nil {
		return err
	}
	destination, err := p.GetAccountByUUID(transaction.DestinationUUID)
	if err != nil {
		return err
	}
	gormTransaction := GormTransaction{
		UUID:            transaction.UUID,
		Status:          transaction.Status,
		SourceUUID:      source.UUID,
		DestinationUUID: destination.UUID,
		Amount:          transaction.Amount,
	}
	err = p.DB.Create(&gormTransaction).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *TestRepo) CreateTransaction(transaction models.Transaction) error {
	_, ok := t.Transaction[transaction.UUID]
	if !ok {
		t.Transaction[transaction.UUID] = &transaction
		sourse, err := t.GetAccountByUUID(transaction.SourceUUID)
		if err != nil {
			return err
		}
		destination, err := t.GetAccountByUUID(transaction.DestinationUUID)
		if err != nil {
			return err
		}
		t.Accounts[sourse.UUID].Sources = append(t.Accounts[sourse.UUID].Sources, transaction)
		t.Accounts[destination.UUID].Destinations = append(t.Accounts[destination.UUID].Destinations, transaction)
		return nil
	}
	return ErrorCreated
}

func (p *PostgresRepo) fromGormToModelAccount(accounts []GormAccount) []models.Account {
	modelAccounts := make([]models.Account, len(accounts))
	for i, acc := range accounts {
		modelAccounts[i] = models.Account{
			UUID:         acc.UUID,
			IBAN:         acc.IBAN,
			Balance:      acc.Balance,
			UserUUID:     acc.UserUUID,
			Sources:      p.fromGormToModelTransaction(acc.Sources),
			Destinations: p.fromGormToModelTransaction(acc.Destinations),
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
func NewTestRepo() TestRepo {
	users := make(map[uuid.UUID]*models.User)
	accounts := make(map[uuid.UUID]*models.Account)
	return TestRepo{
		Users:    users,
		Accounts: accounts,
	}
}

func (p *PostgresRepo) GetAccountByUUID(uuid uuid.UUID) (*models.Account, error) {
	gormAccount := GormAccount{}
	err := p.DB.Model(GormAccount{}).Where("UUID = ?", uuid).Preload("Sources").Preload("Destinations").Take(&gormAccount)
	if err != nil {
		return &models.Account{}, nil
	}
	account := models.Account{
		UUID:         gormAccount.UUID,
		IBAN:         gormAccount.IBAN,
		Balance:      gormAccount.Balance,
		UserUUID:     gormAccount.UserUUID,
		Sources:      p.fromGormToModelTransaction(gormAccount.Sources),
		Destinations: p.fromGormToModelTransaction(gormAccount.Destinations),
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

func (t *TestRepo) GetAccounts(userUUID uuid.UUID) ([]models.Account, error) {
	user, err := t.GetUserByUUID(userUUID)
	if err != nil {
		return user.Accounts, err
	}
	return user.Accounts, nil
}

func (p *TestRepo) GetUserByUUID(uuid uuid.UUID) (*models.User, error) {
	user, ok := p.Users[uuid]
	if !ok {
		return &models.User{}, ErrorUnknownUser
	}
	return user, nil
}

func (p *TestRepo) GetAccountByUUID(uuid uuid.UUID) (*models.Account, error) {
	account, ok := p.Accounts[uuid]
	if !ok {
		return &models.Account{}, ErrorUnknownAccount
	}
	return account, nil
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
