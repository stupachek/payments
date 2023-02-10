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
	GetUser(ctx *gin.Context, email string) models.User
}

type UserPostgresRepo struct {
	gorm.Model
	UUID      uuid.UUID `json:"uuid" gorm:"type:uuid"`
	FisrtName string    `json:"firstName" gorm:"size:50;not null"`
	LastName  string    `json:"lastName" gorm:"size:50;not null"`
	Email     string    `json:"email" gorm:"size:255;not null;unique"`
	Password  string    `json:"password" gorm:"size:250;not null"`
}

type UserTestRepo struct {
	users map[int]models.User
}

func (u *UserTestRepo) CreateUser(ctx *gin.Context, user models.User) error {
	_, ok := u.users[user.ID]
	if !ok {
		u.users[user.ID] = user
		return nil
	}
	ctx.Abort()
	return errors.New("user has already created")
}

func (u *UserPostgresRepo) CreateUser(ctx *gin.Context, user models.User) error {
	err := DB.Create(&u).Error
	if err != nil {
		return err
	}
	return nil
}
