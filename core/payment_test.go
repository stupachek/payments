package core

import (
	"errors"
	"pay/models"
	"pay/repository"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name   string
		user   models.User
		expErr error
	}{
		{
			name: "success",
			user: models.User{FisrtName: "Bob",
				LastName: "Black",
				Email:    "bob.black@gmail.com",
				Password: "bob123"},
			expErr: nil,
		},
	}
	users := make(map[uuid.UUID]models.User)
	userRepo := repository.NewTestRepo(users)
	system := NewPaymentSystem(&userRepo)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := system.Register(&tt.user); !assert.IsEqual(err, tt.expErr) {
				t.Errorf("register error = %v, expErr = %v", err, tt.expErr)
			}
		})
	}

}

func TestRegisterFailed(t *testing.T) {
	users := make(map[uuid.UUID]models.User)
	userRepo := repository.NewTestRepo(users)
	system := NewPaymentSystem(&userRepo)
	user1 := models.User{
		FisrtName: "Bob",
		LastName:  "Black",
		Email:     "bob.black@gmail.com",
		Password:  "bob123",
	}
	user2 := models.User{
		FisrtName: "Bob",
		LastName:  "Right",
		Email:     "bob.black@gmail.com",
		Password:  "bob123",
	}
	if err := system.Register(&user1); err != nil {
		t.Errorf("register error = %v, expErr = %v", err, nil)
	}
	expErr := errors.New("user has already created")
	if err := system.Register(&user2); !assert.IsEqual(err, expErr) {
		t.Errorf("register error = %v, expErr = %v", err, expErr)
	}

}

func TestLogin(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		expErr   error
	}{
		{
			name:     "success",
			email:    "bob.black@gmail.com",
			password: "bob123",
			expErr:   nil,
		},
		{
			name:     "unknown user",
			email:    "alice.go@gmail.com",
			password: "alice123",
			expErr:   ErrUnauthenticated,
		},
		{
			name:     "wrong password",
			email:    "bob.black@gmail.com",
			password: "alice",
			expErr:   ErrUnauthenticated,
		},
	}
	users := make(map[uuid.UUID]models.User)
	userRepo := repository.NewTestRepo(users)
	system := NewPaymentSystem(&userRepo)
	system.Register(&models.User{
		FisrtName: "Bob",
		LastName:  "Right",
		Email:     "bob.black@gmail.com",
		Password:  "bob123",
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := system.LoginCheck(tt.email, tt.password); !assert.IsEqual(err, tt.expErr) {
				t.Errorf("register error = %v, expErr = %v", err, tt.expErr)
			}
		})
	}

}
