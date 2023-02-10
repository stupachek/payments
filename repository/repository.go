package repository

import (
	"errors"
	"pay/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type UserRepository interface {
	CreateUser(ctx *gin.Context, user *models.User) error
	GetUserEmail(ctx *gin.Context, email string) (models.User, error)
	GetUserUUID(ctx *gin.Context, uuid uuid.UUID) (models.User, error)
}

type UserPostgresRepo struct {
	DB *gorm.DB
}

type UserTestRepo struct {
	Users map[uuid.UUID]models.User
}

func (u *UserPostgresRepo) GetUserUUID(ctx *gin.Context, uuid uuid.UUID) (models.User, error) {
	userGorm := models.GormUser{}
	err := u.DB.Model(models.GormUser{}).Where("UUID = ?", uuid).Take(&userGorm).Error
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

func (u *UserTestRepo) GetUserUUID(ctx *gin.Context, uuid uuid.UUID) (models.User, error) {
	user, ok := u.Users[uuid]
	if !ok {
		return models.User{}, errors.New("user does not exist")
	}
	return user, nil
}

func (u *UserPostgresRepo) GetUserEmail(ctx *gin.Context, email string) (models.User, error) {
	userGorm := models.GormUser{}
	err := u.DB.Model(models.GormUser{}).Where("email = ?", email).Take(&userGorm).Error
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

func (u *UserTestRepo) GetUserEmail(ctx *gin.Context, email string) (models.User, error) {
	for _, user := range u.Users {
		if user.Email == "email" {
			return user, nil
		}
	}
	return models.User{}, errors.New("user does not exist")
}

func (u *UserTestRepo) CreateUser(ctx *gin.Context, user *models.User) error {
	_, ok := u.Users[user.UUID]
	if !ok {
		u.Users[user.UUID] = *user
		return nil
	}
	return errors.New("user has already created")
}

func (u *UserPostgresRepo) CreateUser(ctx *gin.Context, user *models.User) error {

	gormU := models.GormUser{
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
