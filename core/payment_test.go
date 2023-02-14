package core

import (
	"errors"
	"pay/models"
	"pay/repository"
	"testing"

	"github.com/go-playground/assert/v2"
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
	testRepo := repository.NewTestRepo()
	system := NewPaymentSystem(&testRepo)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := system.Register(&tt.user); !assert.IsEqual(err, tt.expErr) {
				t.Errorf("register error = %v, expErr = %v", err, tt.expErr)
			}
		})
	}

}

func TestRegisterFailed(t *testing.T) {
	testRepo := repository.NewTestRepo()
	system := NewPaymentSystem(&testRepo)
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
	testRepo := repository.NewTestRepo()
	system := NewPaymentSystem(&testRepo)
	system.Register(&models.User{
		FisrtName: "Bob",
		LastName:  "Black",
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

func TestTokenSuccess(t *testing.T) {
	testRepo := repository.NewTestRepo()
	system := NewPaymentSystem(&testRepo)
	bob := &models.User{
		FisrtName: "Bob",
		LastName:  "Black",
		Email:     "bob.black@gmail.com",
		Password:  "bob123",
	}
	if err := system.Register(bob); err != nil {
		t.Errorf("register error: %v", err)
	}
	token, err := system.LoginCheck("bob.black@gmail.com", "bob123")
	if err != nil {
		t.Errorf("login error: %v", err)
	}
	if err = system.CheckToken(bob.UUID, token); err != nil {
		t.Errorf("token error: %v", err)
	}
}

func TestTokenWrongToken(t *testing.T) {
	testRepo := repository.NewTestRepo()
	system := NewPaymentSystem(&testRepo)
	bob := &models.User{
		FisrtName: "Bob",
		LastName:  "Black",
		Email:     "bob.black@gmail.com",
		Password:  "bob123",
	}
	if err := system.Register(bob); err != nil {
		t.Errorf("register error: %v", err)
	}
	if _, err := system.LoginCheck("bob.black@gmail.com", "bob123"); err != nil {
		t.Errorf("login error: %v", err)
	}
	wrongToken := "c521d0ac2fbea2c9970ac267b5052c55ffba8a2280337bc67334ee218a927d78"
	if err := system.CheckToken(bob.UUID, wrongToken); !assert.IsEqual(err, ErrUnauthenticated) {
		t.Errorf("token error: %v", err)
	}
}

func TestTokenWrongUser(t *testing.T) {
	testRepo := repository.NewTestRepo()
	system := NewPaymentSystem(&testRepo)
	bob := &models.User{
		FisrtName: "Bob",
		LastName:  "Black",
		Email:     "bob.black@gmail.com",
		Password:  "bob123",
	}
	alice := &models.User{
		FisrtName: "Alice",
		LastName:  "Black",
		Email:     "alice.black@gmail.com",
		Password:  "alice123",
	}
	if err := system.Register(bob); err != nil {
		t.Errorf("register error: %v", err)
	}
	if err := system.Register(alice); err != nil {
		t.Errorf("register error: %v", err)
	}
	_, err := system.LoginCheck("bob.black@gmail.com", "bob123")
	if err != nil {
		t.Errorf("login error: %v", err)
	}
	tokenAlice, err := system.LoginCheck("alice.black@gmail.com", "alice123")
	if err != nil {
		t.Errorf("login error: %v", err)
	}
	if err := system.CheckToken(bob.UUID, tokenAlice); !assert.IsEqual(err, ErrUnauthenticated) {
		t.Errorf("token error: %v", err)
	}
}

func TestCreateNewAccountSucces(t *testing.T) {
	testRepo := repository.NewTestRepo()
	system := NewPaymentSystem(&testRepo)
	bob := &models.User{
		FisrtName: "Bob",
		LastName:  "Black",
		Email:     "bob.black@gmail.com",
		Password:  "bob123",
	}
	if err := system.Register(bob); err != nil {
		t.Errorf("register error: %v", err)
	}
	_, err := system.LoginCheck("bob.black@gmail.com", "bob123")
	if err != nil {
		t.Errorf("login error: %v", err)
	}
	if _, err := system.NewAccount(bob.UUID); err != nil {
		t.Errorf("create new account error: %v", err)
	}
}

func TestCreateNewAccountUnknownUser(t *testing.T) {
	testRepo := repository.NewTestRepo()
	system := NewPaymentSystem(&testRepo)
	bob := &models.User{
		FisrtName: "Bob",
		LastName:  "Black",
		Email:     "bob.black@gmail.com",
		Password:  "bob123",
	}
	account, err := system.NewAccount(bob.UUID)
	if !assert.IsEqual(err, repository.ErrorUnknownUser) {
		t.Errorf("create new account error: %v", err)
	}
	if account.UserUUID != bob.UUID {
		t.Errorf("different userID : %v, exp: %v", account.UserUUID, bob.UUID)
	}

}

func TestGetAccounts(t *testing.T) {
	testRepo := repository.NewTestRepo()
	system := NewPaymentSystem(&testRepo)
	bob := &models.User{
		FisrtName: "Bob",
		LastName:  "Black",
		Email:     "bob.black@gmail.com",
		Password:  "bob123",
	}
	if err := system.Register(bob); err != nil {
		t.Errorf("register error: %v", err)
	}
	_, err := system.LoginCheck("bob.black@gmail.com", "bob123")
	if err != nil {
		t.Errorf("login error: %v", err)
	}
	if _, err := system.NewAccount(bob.UUID); err != nil {
		t.Errorf("create new account error: %v", err)
	}
	accs, err := system.GetAccounts(bob.UUID)
	if err != nil {
		t.Errorf("create new account error: %v", err)
	}
	if len(accs) != 1 {
		t.Errorf("diff amount of accounts: %v exp: %v", len(accs), 1)
	}

	if accs[0].UserUUID != bob.UUID {
		t.Errorf("different userUUID : %v, exp: %v", accs[0].UserUUID, bob.UUID)
	}
	if accs[0].Balance != 0 {
		t.Errorf("balance has to be 0")
	}

}

func TestCreateTransaction(t *testing.T) {
	testRepo := repository.NewTestRepo()
	system := NewPaymentSystem(&testRepo)
	bob := &models.User{
		FisrtName: "Bob",
		LastName:  "Black",
		Email:     "bob.black@gmail.com",
		Password:  "bob123",
	}
	if err := system.Register(bob); err != nil {
		t.Errorf("register error: %v", err)
	}
	_, err := system.LoginCheck("bob.black@gmail.com", "bob123")
	if err != nil {
		t.Errorf("login error: %v", err)
	}
	source, err := system.NewAccount(bob.UUID)
	if err != nil {
		t.Errorf("create new account error: %v", err)
	}
	destination, err := system.NewAccount(bob.UUID)
	if err != nil {
		t.Errorf("create new account error: %v", err)
	}
	tr := Transaction{
		UserUUID:        bob.UUID,
		SourceUUID:      source.UUID,
		DestinationUUID: destination.UUID,
		Amount:          0,
	}
	transaction, err := system.NewTransaction(tr)
	if err != nil {
		t.Errorf("create new transaction error: %v", err)
	}
	if transaction.SourceUUID != source.UUID {
		t.Errorf("diff source uuid")
	}
	if transaction.DestinationUUID != destination.UUID {
		t.Errorf("diff destination uuid")
	}
	accs, err := system.GetAccounts(bob.UUID)
	if err != nil {
		t.Errorf("create new account error: %v", err)
	}
	for _, acc := range accs {
		if acc.UUID == tr.SourceUUID {
			if len(acc.Sources) != 1 {
				t.Error("diff sources")
			}
		} else {
			if len(acc.Destinations) != 1 {
				t.Error("diff destinations")
			}
		}
	}

}
