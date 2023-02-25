package integration

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"payment/app"
	"payment/controllers"
	"payment/core"
	"payment/repository"
	"sync"
	"testing"
	"time"

	"github.com/goccy/go-json"
)

func TestPaymentIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	DB := repository.ConnectDataBase()
	repository.ClearData(DB)
	userRepo := repository.NewGormUserRepo(DB)
	system := core.NewPaymentSystem(userRepo)
	controller := controllers.NewHttpController(system)
	err := controller.System.SetupAdmin()
	if err != nil {
		log.Fatalf("can't create admin, err %v", err.Error())
	}
	app := app.New(controller)
	go func() {
		// service connections
		if err := app.Run(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	defer app.Stop()
	time.Sleep(time.Second)
	t.Run("transactionSend", func(t *testing.T) {
		input := controllers.RegisterInput{
			FisrtName: "Bob",
			LastName:  "Blavk",
			Email:     "asdf@qwe.io",
			Password:  "qwerty",
		}
		reqResult := sendReq(t, "POST", "http://localhost:8080/users/register", input, nil)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("error register")
		}
		var userUUID string = reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", "http://localhost:8080/users/login", input, nil)
		token, ok := reqResult["token"].(string)
		if !ok {
			t.Fatal("error login")
		}
		url := fmt.Sprintf("http://localhost:8080/users/%v/accounts/new", userUUID)
		auth := make(map[string]string)
		auth["Authorization"] = token
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		sourceUUID := reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		destinationUUID := reqResult["uuid"].(string)
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/add-money", userUUID, sourceUUID)
		inputAddMoney := controllers.AddMoneyInput{
			Amount: "60",
		}
		reqResult = sendReq(t, "POST", url, inputAddMoney, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("add money error")
		}
		account := reqResult["account"].(map[string]any)
		money := account["balance"].(float64)
		if money != 60 {
			t.Fatalf("balance %v, exp: %v", money, 60)
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/transactions/new", userUUID, sourceUUID)
		inputTr := controllers.TransactionInput{
			DestinationUUID: destinationUUID,
			Amount:          "50",
		}
		reqResult = sendReq(t, "POST", url, inputTr, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create transaction error")
		}
		trans := reqResult["transaction"].(map[string]any)
		transUUID := trans["uuid"].(string)
		log.Print(trans)
		log.Print(transUUID)
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/transactions/%v/send", userUUID, sourceUUID, transUUID)
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("send transaction error")
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v", userUUID, sourceUUID)
		reqResult = sendReq(t, "GET", url, nil, auth)
		if balance := reqResult["balance"].(float64); balance != 10 {
			t.Fatalf("wrong balance :%v, exp:%v", balance, 10)
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v", userUUID, destinationUUID)
		reqResult = sendReq(t, "GET", url, nil, auth)
		if balance := reqResult["balance"].(float64); balance != 50 {
			t.Fatalf("wrong balance :%v, exp:%v", balance, 50)
		}

	})
	t.Run("wrongAccount", func(t *testing.T) {
		inputBob := controllers.RegisterInput{
			FisrtName: "Bob",
			LastName:  "Evans",
			Email:     "evans@gmail.com",
			Password:  "qwerty",
		}
		reqResult := sendReq(t, "POST", "http://localhost:8080/users/register", inputBob, nil)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("error register")
		}
		var userUUIDBob string = reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", "http://localhost:8080/users/login", inputBob, nil)
		tokenBob, ok := reqResult["token"].(string)
		if !ok {
			t.Fatal("error login")
		}
		url := fmt.Sprintf("http://localhost:8080/users/%v/accounts/new", userUUIDBob)
		authBob := make(map[string]string)
		authBob["Authorization"] = tokenBob
		reqResult = sendReq(t, "POST", url, nil, authBob)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		sourceUUID := reqResult["uuid"].(string)
		inputAddMoney := controllers.AddMoneyInput{
			Amount: "60",
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/add-money", userUUIDBob, sourceUUID)
		reqResult = sendReq(t, "POST", url, inputAddMoney, authBob)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("add money error")
		}
		inputAlice := controllers.RegisterInput{
			FisrtName: "Alice",
			LastName:  "Potter",
			Email:     "potter@gmail.com",
			Password:  "potter123",
		}
		reqResult = sendReq(t, "POST", "http://localhost:8080/users/register", inputAlice, nil)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("error register")
		}
		var userUUIDAlice string = reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", "http://localhost:8080/users/login", inputAlice, nil)
		tokenAlice, ok := reqResult["token"].(string)
		if !ok {
			t.Fatal("error login")
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/new", userUUIDAlice)
		authAlice := make(map[string]string)
		authAlice["Authorization"] = tokenAlice
		reqResult = sendReq(t, "POST", url, nil, authAlice)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		destinationUUID := reqResult["uuid"].(string)
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/transactions/new", userUUIDAlice, sourceUUID)
		inputTr := controllers.TransactionInput{
			DestinationUUID: destinationUUID,
			Amount:          "50",
		}
		reqResult = sendReq(t, "POST", url, inputTr, authAlice)
		if err := reqResult["error"]; err != "wrong account" {
			t.Fatalf("error: %v, exp:%v", err, "wrong account")
		}

	})
	t.Run("unauthenticated", func(t *testing.T) {
		inputBob := controllers.RegisterInput{
			FisrtName: "Bob",
			LastName:  "Potter",
			Email:     "bob.potter@gmail.com",
			Password:  "qwerty",
		}
		reqResult := sendReq(t, "POST", "http://localhost:8080/users/register", inputBob, nil)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("error register")
		}
		var userUUIDBob string = reqResult["uuid"].(string)
		url := fmt.Sprintf("http://localhost:8080/users/%v/accounts/new", userUUIDBob)
		reqResult = sendReq(t, "POST", url, nil, nil)
		if err := reqResult["error"]; err != "unauthenticated" {
			t.Fatalf("error: %v, exp: %v", err, "unauthenticated")
		}

	})
	t.Run("insufficientFunds", func(t *testing.T) {
		input := controllers.RegisterInput{
			FisrtName: "Bob",
			LastName:  "Smith",
			Email:     "smith@ui.o",
			Password:  "qwerty",
		}
		reqResult := sendReq(t, "POST", "http://localhost:8080/users/register", input, nil)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("error register")
		}
		var userUUID string = reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", "http://localhost:8080/users/login", input, nil)
		token, ok := reqResult["token"].(string)
		if !ok {
			t.Fatal("error login")
		}
		url := fmt.Sprintf("http://localhost:8080/users/%v/accounts/new", userUUID)
		auth := make(map[string]string)
		auth["Authorization"] = token
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		sourceUUID := reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		destinationUUID := reqResult["uuid"].(string)
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/add-money", userUUID, sourceUUID)
		inputAddMoney := controllers.AddMoneyInput{
			Amount: "60",
		}
		reqResult = sendReq(t, "POST", url, inputAddMoney, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("add money error")
		}
		account := reqResult["account"].(map[string]any)
		money := account["balance"].(float64)
		if money != 60 {
			t.Fatalf("balance %v, exp: %v", money, 60)
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/transactions/new", userUUID, sourceUUID)
		inputTr := controllers.TransactionInput{
			DestinationUUID: destinationUUID,
			Amount:          "100",
		}
		reqResult = sendReq(t, "POST", url, inputTr, auth)
		if err := reqResult["error"]; err != "insufficient funds" {
			t.Fatalf("error: %v, exp: %v", err, "insufficient funds")
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v", userUUID, sourceUUID)
		reqResult = sendReq(t, "GET", url, nil, auth)
		if balance := reqResult["balance"].(float64); balance != 60 {
			t.Fatalf("wrong balance :%v, exp:%v", balance, 60)
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v", userUUID, destinationUUID)
		reqResult = sendReq(t, "GET", url, nil, auth)
		if balance := reqResult["balance"].(float64); balance != 0 {
			t.Fatalf("wrong balance :%v, exp:%v", balance, 0)
		}

	})
	t.Run("insufficientFundsOnSend", func(t *testing.T) {
		input := controllers.RegisterInput{
			FisrtName: "Bob",
			LastName:  "Taylor",
			Email:     "taylor.f@i.ua",
			Password:  "qwerty",
		}
		reqResult := sendReq(t, "POST", "http://localhost:8080/users/register", input, nil)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("error register")
		}
		var userUUID string = reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", "http://localhost:8080/users/login", input, nil)
		token, ok := reqResult["token"].(string)
		if !ok {
			t.Fatal("error login")
		}
		url := fmt.Sprintf("http://localhost:8080/users/%v/accounts/new", userUUID)
		auth := make(map[string]string)
		auth["Authorization"] = token
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		sourceUUID := reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		destinationUUID := reqResult["uuid"].(string)
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/add-money", userUUID, sourceUUID)
		inputAddMoney := controllers.AddMoneyInput{
			Amount: "100",
		}
		reqResult = sendReq(t, "POST", url, inputAddMoney, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("add money error")
		}
		account := reqResult["account"].(map[string]any)
		money := account["balance"].(float64)
		if money != 100 {
			t.Fatalf("balance %v, exp: %v", money, 100)
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/transactions/new", userUUID, sourceUUID)
		inputTr1 := controllers.TransactionInput{
			DestinationUUID: destinationUUID,
			Amount:          "70",
		}
		reqResult = sendReq(t, "POST", url, inputTr1, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatalf("create transaction error")
		}
		trans1 := reqResult["transaction"].(map[string]any)
		trans1UUID := trans1["uuid"].(string)

		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/transactions/new", userUUID, sourceUUID)
		inputTr2 := controllers.TransactionInput{
			DestinationUUID: destinationUUID,
			Amount:          "50",
		}
		reqResult = sendReq(t, "POST", url, inputTr2, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatalf("create transaction error")
		}
		trans2 := reqResult["transaction"].(map[string]any)
		trans2UUID := trans2["uuid"].(string)

		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/transactions/%v/send", userUUID, sourceUUID, trans1UUID)
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("send transaction error")
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/transactions/%v/send", userUUID, sourceUUID, trans2UUID)
		reqResult = sendReq(t, "POST", url, nil, auth)
		if err := reqResult["error"]; err != "insufficient funds" {
			t.Fatalf("error: %v, exp: %v", err, "insufficient funds")
		}

		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v", userUUID, sourceUUID)
		reqResult = sendReq(t, "GET", url, nil, auth)
		if balance := reqResult["balance"].(float64); balance != 30 {
			t.Fatalf("wrong balance :%v, exp:%v", balance, 30)
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v", userUUID, destinationUUID)
		reqResult = sendReq(t, "GET", url, nil, auth)
		if balance := reqResult["balance"].(float64); balance != 70 {
			t.Fatalf("wrong balance :%v, exp:%v", balance, 70)
		}

	})
	t.Run("transaction", func(t *testing.T) {
		input := controllers.RegisterInput{
			FisrtName: "Bob",
			LastName:  "Lobster",
			Email:     "lobster@i.ua",
			Password:  "qwerty",
		}
		reqResult := sendReq(t, "POST", "http://localhost:8080/users/register", input, nil)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("error register")
		}
		var userUUID string = reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", "http://localhost:8080/users/login", input, nil)
		token, ok := reqResult["token"].(string)
		if !ok {
			t.Fatal("error login")
		}
		url := fmt.Sprintf("http://localhost:8080/users/%v/accounts/new", userUUID)
		auth := make(map[string]string)
		auth["Authorization"] = token
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		sourceUUID := reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		destinationUUID := reqResult["uuid"].(string)
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/add-money", userUUID, sourceUUID)
		inputAddMoney := controllers.AddMoneyInput{
			Amount: "100",
		}
		reqResult = sendReq(t, "POST", url, inputAddMoney, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("add money error")
		}
		account := reqResult["account"].(map[string]any)
		money := account["balance"].(float64)
		if money != 100 {
			t.Fatalf("balance %v, exp: %v", money, 100)
		}
		var wg sync.WaitGroup
		wg.Add(50)
		for i := 0; i < 50; i++ {
			go func() {
				defer wg.Done()
				url := fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/transactions/new", userUUID, sourceUUID)
				inputTr := controllers.TransactionInput{
					DestinationUUID: destinationUUID,
					Amount:          "1",
				}
				reqResult := sendReq(t, "POST", url, inputTr, auth)
				if _, ok := reqResult["message"]; !ok {
					log.Fatal("create transaction error")
				}
				trans := reqResult["transaction"].(map[string]any)
				transUUID := trans["uuid"].(string)
				url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/transactions/%v/send", userUUID, sourceUUID, transUUID)
				reqResult = sendReq(t, "POST", url, nil, auth)
				if _, ok := reqResult["message"]; !ok {
					log.Fatal("send transaction error")
				}

			}()
		}
		wg.Wait()
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v", userUUID, destinationUUID)
		reqResult = sendReq(t, "GET", url, nil, auth)
		if balance := reqResult["balance"].(float64); balance != 50 {
			t.Fatalf("wrong balance :%v, exp:%v", balance, 50)
		}

	})
	t.Run("queryAccountsSuccess", func(t *testing.T) {
		input := controllers.RegisterInput{
			FisrtName: "Bob",
			LastName:  "Williams",
			Email:     "bob.www@qmail.com",
			Password:  "qwerty",
		}
		reqResult := sendReq(t, "POST", "http://localhost:8080/users/register", input, nil)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("error register")
		}
		var userUUID string = reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", "http://localhost:8080/users/login", input, nil)
		token, ok := reqResult["token"].(string)
		if !ok {
			t.Fatal("error login")
		}
		url := fmt.Sprintf("http://localhost:8080/users/%v/accounts/new", userUUID)
		auth := make(map[string]string)
		auth["Authorization"] = token
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts?sort_by=iban&orser=desc&limit=3&offset=1", userUUID)
		reqResult = sendReq(t, "GET", url, nil, auth)
		acc, ok := reqResult["accounts"].([]any)
		if !ok {
			t.Fatal("get accounts error")
		}
		if len(acc) != 3 {
			t.Fatalf("len: %v, exp: %v", len(acc), 3)
		}

	})
	t.Run("queryAccountssuccessFailed", func(t *testing.T) {
		input := controllers.RegisterInput{
			FisrtName: "Bob",
			LastName:  "Stevens",
			Email:     "bob.stevens@qmail.com",
			Password:  "qwerty",
		}
		reqResult := sendReq(t, "POST", "http://localhost:8080/users/register", input, nil)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("error register")
		}
		var userUUID string = reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", "http://localhost:8080/users/login", input, nil)
		token, ok := reqResult["token"].(string)
		if !ok {
			t.Fatal("error login")
		}
		url := fmt.Sprintf("http://localhost:8080/users/%v/accounts/new", userUUID)
		auth := make(map[string]string)
		auth["Authorization"] = token
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}

		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts?sort_by=qwerty&orser=desc&limit=3&offset=1", userUUID)
		reqResult = sendReq(t, "GET", url, nil, auth)
		err := reqResult["error"].(string)
		if err != "unknown query" {
			t.Fatalf("get accounts error: %v, exp: %v", err, "unknown query")
		}

	})
	t.Run("queryTransactionsSucces", func(t *testing.T) {
		input := controllers.RegisterInput{
			FisrtName: "Bob",
			LastName:  "Fox",
			Email:     "fox.fff@gmail.com",
			Password:  "qwerty",
		}
		reqResult := sendReq(t, "POST", "http://localhost:8080/users/register", input, nil)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("error register")
		}
		var userUUID string = reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", "http://localhost:8080/users/login", input, nil)
		token, ok := reqResult["token"].(string)
		if !ok {
			t.Fatal("error login")
		}
		url := fmt.Sprintf("http://localhost:8080/users/%v/accounts/new", userUUID)
		auth := make(map[string]string)
		auth["Authorization"] = token
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		sourceUUID := reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		destinationUUID := reqResult["uuid"].(string)
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/add-money", userUUID, sourceUUID)
		inputAddMoney := controllers.AddMoneyInput{
			Amount: "60",
		}
		reqResult = sendReq(t, "POST", url, inputAddMoney, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("add money error")
		}
		account := reqResult["account"].(map[string]any)
		money := account["balance"].(float64)
		if money != 60 {
			t.Fatalf("balance %v, exp: %v", money, 60)
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/transactions/new", userUUID, sourceUUID)
		inputTr := controllers.TransactionInput{
			DestinationUUID: destinationUUID,
			Amount:          "1",
		}
		reqResult = sendReq(t, "POST", url, inputTr, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create transaction error")
		}
		reqResult = sendReq(t, "POST", url, inputTr, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create transaction error")
		}
		reqResult = sendReq(t, "POST", url, inputTr, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create transaction error")
		}
		reqResult = sendReq(t, "POST", url, inputTr, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create transaction error")
		}
		reqResult = sendReq(t, "POST", url, inputTr, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create transaction error")
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/transactions?limit=4&offset=1&sort_by=uuid", userUUID, sourceUUID)
		reqResult = sendReq(t, "GET", url, nil, auth)
		tr, ok := reqResult["transactions"].([]any)
		if !ok {
			t.Fatal("get transaction error")
		}
		if len(tr) != 4 {
			t.Fatalf("len: %v, exp: %v", len(tr), 3)
		}

	})
	t.Run("queryTransactionsFailed", func(t *testing.T) {
		input := controllers.RegisterInput{
			FisrtName: "Bob",
			LastName:  "West",
			Email:     "bob.west123@gmail.com",
			Password:  "qwerty",
		}
		reqResult := sendReq(t, "POST", "http://localhost:8080/users/register", input, nil)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("error register")
		}
		var userUUID string = reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", "http://localhost:8080/users/login", input, nil)
		token, ok := reqResult["token"].(string)
		if !ok {
			t.Fatal("error login")
		}
		url := fmt.Sprintf("http://localhost:8080/users/%v/accounts/new", userUUID)
		auth := make(map[string]string)
		auth["Authorization"] = token
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		sourceUUID := reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		destinationUUID := reqResult["uuid"].(string)
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/add-money", userUUID, sourceUUID)
		inputAddMoney := controllers.AddMoneyInput{
			Amount: "60",
		}
		reqResult = sendReq(t, "POST", url, inputAddMoney, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("add money error")
		}
		account := reqResult["account"].(map[string]any)
		money := account["balance"].(float64)
		if money != 60 {
			t.Fatalf("balance %v, exp: %v", money, 60)
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/transactions/new", userUUID, sourceUUID)
		inputTr := controllers.TransactionInput{
			DestinationUUID: destinationUUID,
			Amount:          "1",
		}
		reqResult = sendReq(t, "POST", url, inputTr, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create transaction error")
		}
		reqResult = sendReq(t, "POST", url, inputTr, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create transaction error")
		}
		reqResult = sendReq(t, "POST", url, inputTr, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create transaction error")
		}
		reqResult = sendReq(t, "POST", url, inputTr, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create transaction error")
		}
		reqResult = sendReq(t, "POST", url, inputTr, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create transaction error")
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/transactions?limit=qwerty", userUUID, sourceUUID)
		reqResult = sendReq(t, "GET", url, nil, auth)
		err := reqResult["error"].(string)
		if err != "unknown query" {
			t.Fatalf("get accounts error: %v, exp: %v", err, "unknown query")
		}

	})
	t.Run("newAdmin", func(t *testing.T) {
		inputBob := controllers.RegisterInput{
			FisrtName: "Bob",
			LastName:  "Moss",
			Email:     "bob.moss@gmail.com",
			Password:  "qwerty",
		}
		reqResult := sendReq(t, "POST", "http://localhost:8080/users/register", inputBob, nil)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("error register")
		}
		var UUIDBob string = reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", "http://localhost:8080/users/login", inputBob, nil)
		_, ok := reqResult["token"].(string)
		if !ok {
			t.Fatal("error login")
		}
		inputAdmin := controllers.LoginInput{
			Email:    "admin@admin.admin",
			Password: "admin",
		}
		reqResult = sendReq(t, "POST", "http://localhost:8080/users/login", inputAdmin, nil)
		tokenAdmin, ok := reqResult["token"].(string)
		if !ok {
			t.Fatal("error login")
		}

		url := fmt.Sprintf("http://localhost:8080/admin/%v/update-role", reqResult["uuid"].(string))
		auth := make(map[string]string)
		auth["Authorization"] = tokenAdmin
		inputRole := controllers.ChangeRoleInput{
			UserUUID: UUIDBob,
			Role:     "admin",
		}
		reqResult = sendReq(t, "POST", url, inputRole, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("error change role")
		}

	})
	t.Run("newAdminFailed", func(t *testing.T) {
		inputBob := controllers.RegisterInput{
			FisrtName: "Bob",
			LastName:  "Evans",
			Email:     "bob.evans123@gmail.com",
			Password:  "qwerty",
		}
		reqResult := sendReq(t, "POST", "http://localhost:8080/users/register", inputBob, nil)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("error register")
		}
		var UUIDBob string = reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", "http://localhost:8080/users/login", inputBob, nil)
		_, ok := reqResult["token"].(string)
		if !ok {
			t.Fatal("error login")
		}
		inputAdmin := controllers.LoginInput{
			Email:    "admin@admin.admin",
			Password: "admin",
		}
		reqResult = sendReq(t, "POST", "http://localhost:8080/users/login", inputAdmin, nil)
		tokenAdmin, ok := reqResult["token"].(string)
		if !ok {
			t.Fatal("error login")
		}

		url := fmt.Sprintf("http://localhost:8080/admin/%v/update-role", reqResult["uuid"].(string))
		auth := make(map[string]string)
		auth["Authorization"] = tokenAdmin
		inputRole := controllers.ChangeRoleInput{
			UserUUID: UUIDBob,
			Role:     "superman",
		}
		reqResult = sendReq(t, "POST", url, inputRole, auth)
		if err := reqResult["error"]; err != controllers.UnknownRoleError {
			t.Fatalf("error change role: %v, exp %v", err, controllers.UnknownRoleError)
		}

	})
	t.Run("blockUnblock", func(t *testing.T) {
		inputBob := controllers.RegisterInput{
			FisrtName: "Bob",
			LastName:  "Lee",
			Email:     "lee.bee@gmail.com",
			Password:  "qwerty",
		}
		reqResult := sendReq(t, "POST", "http://localhost:8080/users/register", inputBob, nil)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("error register")
		}
		var UUIDBob string = reqResult["uuid"].(string)
		reqResult = sendReq(t, "POST", "http://localhost:8080/users/login", inputBob, nil)
		tokenBob, ok := reqResult["token"].(string)
		if !ok {
			t.Fatal("error login")
		}
		url := fmt.Sprintf("http://localhost:8080/users/%v/accounts/new", UUIDBob)
		auth := make(map[string]string)
		auth["Authorization"] = tokenBob
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		accountUUID := reqResult["uuid"].(string)
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/block", UUIDBob, accountUUID)
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("block account error")
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/unblock", UUIDBob, accountUUID)
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("request to unblock account error")
		}
		inputAdmin := controllers.LoginInput{
			Email:    "admin@admin.admin",
			Password: "admin",
		}
		reqResult = sendReq(t, "POST", "http://localhost:8080/users/login", inputAdmin, nil)
		tokenAdmin, ok := reqResult["token"].(string)
		if !ok {
			t.Fatal("error login")
		}
		adminUUID := reqResult["uuid"].(string)
		url = fmt.Sprintf("http://localhost:8080/admin/%v/accounts/%v/unblock", adminUUID, accountUUID)
		auth["Authorization"] = tokenAdmin
		reqResult = sendReq(t, "POST", url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("unblock account error")
		}

	})

}

func sendReq(t *testing.T, method string, url string, inputT interface{}, headers map[string]string) map[string]any {
	input, err := json.Marshal(inputT)
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(input))
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Content-Type", "application/json")
	for name, value := range headers {
		req.Header.Add(name, value)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	var reqResult map[string]any
	err = json.NewDecoder(res.Body).Decode(&reqResult)
	if err != nil {
		t.Error(err)
	}
	return reqResult

}
