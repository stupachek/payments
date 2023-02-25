package repository

import (
	"payment/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUUID(uuid uuid.UUID) (*models.User, error)
	CreateAccount(account *models.Account) error
	CreateTransaction(transaction models.Transaction) error
	GetAccountsForUser(userUUID uuid.UUID, query models.QueryParams) ([]models.Account, error)
	GetAccountByUUID(uuid uuid.UUID) (*models.Account, error)
	GetTransactionForAccount(accountUUID uuid.UUID, query models.QueryParams) ([]models.Transaction, error)
	GetTransactionByUUID(transactionUUID uuid.UUID) (*models.Transaction, error)
	IncBalance(accountUUID uuid.UUID, amount uint) error
	DecBalance(accountUUID uuid.UUID, amount uint) error
	UpdateStatus(transactionUUID uuid.UUID, status string) error
	Transaction(callback func(repo Repository) error) error
	UpdateRole(userUUID uuid.UUID, role string) error
	UpdatePassword(userUUID uuid.UUID, password string) error
}

type PostgresRepo struct {
	DB *gorm.DB
}

func (p *PostgresRepo) Transaction(callback func(repo Repository) error) error {
	return p.DB.Transaction(func(tx *gorm.DB) error {
		repo := PostgresRepo{
			DB: tx,
		}
		return callback(&repo)
	})
}

func (p *PostgresRepo) IncBalance(accountUUID uuid.UUID, amount uint) error {
	return p.DB.Model(&GormAccount{}).Where("UUID = ?", accountUUID).Update("Balance", gorm.Expr("Balance + ?", amount)).Error
}

func (p *PostgresRepo) DecBalance(accountUUID uuid.UUID, amount uint) error {
	return p.DB.Model(&GormAccount{}).Where("UUID = ?", accountUUID).Update("Balance", gorm.Expr("Balance - ?", amount)).Error
}

func (p *PostgresRepo) UpdateStatus(transactionUUID uuid.UUID, status string) error {
	return p.DB.Model(&GormTransaction{}).Where("UUID = ?", transactionUUID).Update("Status", status).Error
}

func (p *PostgresRepo) UpdateRole(userUUID uuid.UUID, role string) error {
	return p.DB.Model(&GormUser{}).Where("UUID = ?", userUUID).Update("Role", role).Error
}
func (p *PostgresRepo) UpdatePassword(userUUID uuid.UUID, password string) error {
	return p.DB.Model(&GormUser{}).Where("UUID = ?", userUUID).Update("Password", password).Error
}

func (p *PostgresRepo) GetTransactionByUUID(transactionUUID uuid.UUID) (*models.Transaction, error) {
	var gormTransaction GormTransaction
	result := p.DB.Model(&GormTransaction{}).Where("UUID = ?", transactionUUID).Take(&gormTransaction)
	if err := result.Error; err != nil {
		return &models.Transaction{}, err
	}
	transaction := models.Transaction{
		UUID:            gormTransaction.UUID,
		Status:          gormTransaction.Status,
		SourceUUID:      gormTransaction.SourceUUID,
		DestinationUUID: gormTransaction.DestinationUUID,
		Amount:          gormTransaction.Amount,
		CreatedAt:       gormTransaction.CreatedAt,
		UpdatedAt:       gormTransaction.UpdatedAt,
	}
	return &transaction, nil
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

func (p *PostgresRepo) GetTransactionForAccount(accountUUID uuid.UUID, query models.QueryParams) ([]models.Transaction, error) {
	var gormTransaction []GormTransaction
	result := p.DB.Model(GormTransaction{}).Where("Source_UUID = ? OR Destination_UUID = ?", accountUUID, accountUUID).Order(query.Sort).Limit(int(query.Limit)).Offset(int(query.Offset)).Find(&gormTransaction)
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
			Status:   acc.Status,
		}
	}

	return modelAccounts
}

func (p *PostgresRepo) fromGormToModelTransaction(transactions []GormTransaction) []models.Transaction {
	modelTransaction := make([]models.Transaction, len(transactions))
	for i, tr := range transactions {
		modelTransaction[i] = models.Transaction{
			UUID:            tr.UUID,
			Status:          tr.Status,
			SourceUUID:      tr.SourceUUID,
			DestinationUUID: tr.DestinationUUID,
			Amount:          tr.Amount,
			CreatedAt:       tr.CreatedAt,
			UpdatedAt:       tr.UpdatedAt,
		}
	}
	return modelTransaction
}

func (p *PostgresRepo) GetAccountsForUser(userUUID uuid.UUID, query models.QueryParams) ([]models.Account, error) {
	var gormAccounts []GormAccount
	result := p.DB.Model(GormAccount{}).Where("User_UUID = ?", userUUID).Order(query.Sort).Limit(int(query.Limit)).Offset(int(query.Offset)).Find(&gormAccounts)
	if err := result.Error; err != nil {
		return []models.Account{}, err
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
		Status:   account.Status,
	}
	err := p.DB.Create(&gormAcc).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepo) GetAccountByUUID(uuid uuid.UUID) (*models.Account, error) {
	gormAccount := GormAccount{}
	err := p.DB.Model(GormAccount{}).Where("UUID = ?", uuid).Take(&gormAccount).Error
	if err != nil {
		return &models.Account{}, nil
	}
	account := models.Account{
		UUID:     gormAccount.UUID,
		IBAN:     gormAccount.IBAN,
		Balance:  gormAccount.Balance,
		UserUUID: gormAccount.UserUUID,
		Status:   gormAccount.Status,
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
		Role:      userGorm.Role,
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
		Role:      userGorm.Role,
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
		Role:      user.Role,
		Accounts:  []GormAccount{},
	}
	err := p.DB.Create(&gormU).Error
	if err != nil {
		return err
	}
	return nil
}

func NewGormUserRepo(DB *gorm.DB) Repository {
	return &PostgresRepo{
		DB: DB,
	}
}
