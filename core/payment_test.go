package core

import (
	"pay/models"
	"pay/repository"
	"testing"

	"github.com/google/uuid"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name   string
		user   models.User
		expErr error
	}{
		{
			name: "succes",
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
			if err := system.Register(&tt.user); err != tt.expErr {
				t.Errorf("register error = %v, expErr = %v", err, tt.expErr)
			}
		})
	}

}
