package repository

import (
	"errors"
	"pay/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type UserRepository interface {
	CreateUser(ctx *gin.Context, user models.User) error
	GetUser(ctx *gin.Context, email string) (models.User, error)
}

type UserPostgresRepo struct {
	DB *gorm.DB
}

type UserTestRepo struct {
	users map[uuid.UUID]models.User
}

func (u *UserPostgresRepo) GetUser(ctx *gin.Context, email string) (models.User, error) {
	user := models.User{}
	err := u.DB.Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *UserTestRepo) GetUser(ctx *gin.Context, email string) (models.User, error) {
	for _, user := range u.users {
		if user.Email == "email" {
			return user, nil
		}
	}
	return models.User{}, errors.New("user does not exist")
}

func (u *UserTestRepo) CreateUser(ctx *gin.Context, user models.User) error {
	_, ok := u.users[user.UUID]
	if !ok {
		u.users[user.UUID] = user
		return nil
	}
	return errors.New("user has already created")
}

func (u *UserPostgresRepo) CreateUser(ctx *gin.Context, user models.User) error {
	type gormUser struct {
		gorm.Model
		UUID      uuid.UUID `json:"uuid" gorm:"type:uuid"`
		FisrtName string    `json:"firstName" gorm:"size:50;not null"`
		LastName  string    `json:"lastName" gorm:"size:50;not null"`
		Email     string    `json:"email" gorm:"size:255;not null;unique"`
		Password  string    `json:"password" gorm:"size:250;not null"`
	}
	gormU := gormUser{
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
