package repository

import (
	"errors"
	"pay/models"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var ErrorCreated = errors.New("user has already created")
var ErrorUnknownUser = errors.New("user does not exist")

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUUID(uuid uuid.UUID) (*models.User, error)
	CreateAccount(account *models.Account) error
	GetAccountsForUserWith(uuid uuid.UUID) ([]models.Account, error)
}

type PostgresRepo struct {
	DB *gorm.DB
}

type TestRepo struct {
	idCounter uint
	Users     map[uuid.UUID]*models.User
	Accounts  map[uuid.UUID]*models.Account
}

// TODO : transaction
func fromGormToModelAccount(accounts []GormAccount) []models.Account {
	modelAccounts := make([]models.Account, len(accounts))
	for i, acc := range accounts {
		modelAccounts[i] = models.Account{
			ID:      acc.ID,
			UUID:    acc.UUID,
			IBAN:    acc.IBAN,
			Balance: acc.Balance,
			UserId:  acc.UserId,
		}
	}

	return modelAccounts
}

func (p *PostgresRepo) CreateAccount(account *models.Account) error {
	gormAcc := GormAccount{
		UUID:    account.UUID,
		IBAN:    account.IBAN,
		Balance: account.Balance,
		UserId:  account.UserId,
	}
	err := p.DB.Create(&gormAcc).Error
	if err != nil {
		return err
	}
	return nil
}
func NewTestRepo(users map[uuid.UUID]*models.User, accounts map[uuid.UUID]*models.Account) TestRepo {
	return TestRepo{
		Users:    users,
		Accounts: accounts,
	}
}

func (p *PostgresRepo) GetUserByUUID(uuid uuid.UUID) (*models.User, error) {
	userGorm := GormUser{}
	err := p.DB.Model(GormUser{}).Where("UUID = ?", uuid).Preload("Accounts").Take(&userGorm).Error
	if err != nil {
		return &models.User{}, err
	}
	user := models.User{
		ID:        userGorm.ID,
		UUID:      userGorm.UUID,
		FisrtName: userGorm.FisrtName,
		LastName:  userGorm.LastName,
		Email:     userGorm.Email,
		Password:  userGorm.Password,
		Accounts:  fromGormToModelAccount(userGorm.Accounts),
	}
	return &user, nil
}

func (p *PostgresRepo) GetAccountsForUserWith(uuid uuid.UUID) ([]models.Account, error) {
	user, err := p.GetUserByUUID(uuid)
	if err != nil {
		return []models.Account{}, err
	}
	return user.Accounts, nil
}

func (t *TestRepo) GetAccountsForUserWith(uuid uuid.UUID) ([]models.Account, error) {
	user, err := t.GetUserByUUID(uuid)
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

func (p *PostgresRepo) GetUserByEmail(email string) (*models.User, error) {
	userGorm := GormUser{}
	err := p.DB.Model(GormUser{}).Where("email = ?", email).Preload("Accounts").Take(&userGorm).Error
	if err != nil {
		return &models.User{}, err
	}
	user := models.User{
		ID:        userGorm.ID,
		UUID:      userGorm.UUID,
		FisrtName: userGorm.FisrtName,
		LastName:  userGorm.LastName,
		Email:     userGorm.Email,
		Password:  userGorm.Password,
		Accounts:  fromGormToModelAccount(userGorm.Accounts),
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
	account.ID = t.nextId()
	_, ok := t.Accounts[account.UUID]
	if !ok {
		t.Accounts[account.UUID] = account
		user := t.getUserById(account.UserId)
		user.Accounts = append(user.Accounts, *account)
		return nil
	}
	return ErrorCreated
}

func (t *TestRepo) getUserById(ID uint) *models.User {
	for _, user := range t.Users {
		if user.ID == ID {
			return user
		}
	}
	return &models.User{}
}

func (t *TestRepo) CreateUser(user *models.User) error {
	user.ID = t.nextId()
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

func (t *TestRepo) nextId() uint {
	t.idCounter++
	return t.idCounter
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
		Model:     gorm.Model{},
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
