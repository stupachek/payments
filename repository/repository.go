package repository

import (
	"errors"
	"pay/models"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (models.User, error)
	GetUserByUUID(uuid uuid.UUID) (models.User, error)
}

type UserPostgresRepo struct {
	DB *gorm.DB
}

type UserTestRepo struct {
	Users map[uuid.UUID]models.User
}

func NewTestRepo(users map[uuid.UUID]models.User) UserTestRepo {
	return UserTestRepo{
		Users: users,
	}
}

func (u *UserPostgresRepo) GetUserByUUID(uuid uuid.UUID) (models.User, error) {
	userGorm := GormUser{}
	err := u.DB.Model(GormUser{}).Where("UUID = ?", uuid).Take(&userGorm).Error
	if err != nil {
		return models.User{}, err
	}
	user := models.User{
		UUID:      userGorm.UUID,
		FisrtName: userGorm.FisrtName,
		LastName:  userGorm.LastName,
		Email:     userGorm.Email,
		Password:  userGorm.Password,
	}
	return user, nil
}

func (u *UserTestRepo) GetUserByUUID(uuid uuid.UUID) (models.User, error) {
	user, ok := u.Users[uuid]
	if !ok {
		return models.User{}, errors.New("user does not exist")
	}
	return user, nil
}

func (u *UserPostgresRepo) GetUserByEmail(email string) (models.User, error) {
	userGorm := GormUser{}
	err := u.DB.Model(GormUser{}).Where("email = ?", email).Take(&userGorm).Error
	if err != nil {
		return models.User{}, err
	}
	user := models.User{
		UUID:      userGorm.UUID,
		FisrtName: userGorm.FisrtName,
		LastName:  userGorm.LastName,
		Email:     userGorm.Email,
		Password:  userGorm.Password,
	}
	return user, nil
}

func (u *UserTestRepo) GetUserByEmail(email string) (models.User, error) {
	for _, user := range u.Users {
		if user.Email == "email" {
			return user, nil
		}
	}
	return models.User{}, errors.New("user does not exist")
}

func (u *UserTestRepo) CreateUser(user *models.User) error {
	_, ok := u.Users[user.UUID]
	if !ok {
		u.Users[user.UUID] = *user
		return nil
	}
	return errors.New("user has already created")
}

func (u *UserPostgresRepo) CreateUser(user *models.User) error {

	gormU := GormUser{
		Model:     gorm.Model{},
		UUID:      user.UUID,
		FisrtName: user.FisrtName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	}
	err := u.DB.Create(&gormU).Error
	if err != nil {
		return err
	}
	return nil
}

func NewGormUserRepo(DB *gorm.DB) UserRepository {
	return &UserPostgresRepo{
		DB: DB,
	}
}
