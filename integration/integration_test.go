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
	"testing"
	"time"

	"github.com/goccy/go-json"
)

func TestPaymentIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	params := repository.DBParams{
		Dbdriver: "postgres",
		Host:     "127.0.0.1 ",
		User:     "postgres",
		Password: "postgres",
		Name:     "qwerty",
		Port:     "5432",
	}
	DB := repository.ConnectDataBaseWithParams(params)
	repository.ClearData(DB)
	userRepo := repository.NewGormUserRepo(DB)
	system := core.NewPaymentSystem(userRepo)
	controller := controllers.NewHttpController(system)
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
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/show-balance", userUUID, sourceUUID)
		reqResult = sendReq(t, "GET", url, nil, auth)
		if balance := reqResult["balance"].(float64); balance != 10 {
			t.Fatalf("wrong balance :%v, exp:%v", balance, 10)
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/show-balance", userUUID, destinationUUID)
		reqResult = sendReq(t, "GET", url, nil, auth)
		if balance := reqResult["balance"].(float64); balance != 50 {
			t.Fatalf("wrong balance :%v, exp:%v", balance, 50)
		}

	})
	t.Run("wrongAccount", func(t *testing.T) {
		inputBob := controllers.RegisterInput{
			FisrtName: "Bob",
			LastName:  "Blavk",
			Email:     "asdf@qwe.io",
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
			LastName:  "Blavk",
			Email:     "asdf@qwe.io",
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
			Amount:          "100",
		}
		reqResult = sendReq(t, "POST", url, inputTr, auth)
		if err := reqResult["error"]; err != "insufficient funds" {
			t.Fatalf("error: %v, exp: %v", err, "insufficient funds")
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/show-balance", userUUID, sourceUUID)
		reqResult = sendReq(t, "GET", url, nil, auth)
		if balance := reqResult["balance"].(float64); balance != 60 {
			t.Fatalf("wrong balance :%v, exp:%v", balance, 60)
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/show-balance", userUUID, destinationUUID)
		reqResult = sendReq(t, "GET", url, nil, auth)
		if balance := reqResult["balance"].(float64); balance != 0 {
			t.Fatalf("wrong balance :%v, exp:%v", balance, 0)
		}

	})
	t.Run("insufficientFundsOnSend", func(t *testing.T) {
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

		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/show-balance", userUUID, sourceUUID)
		reqResult = sendReq(t, "GET", url, nil, auth)
		if balance := reqResult["balance"].(float64); balance != 30 {
			t.Fatalf("wrong balance :%v, exp:%v", balance, 30)
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/show-balance", userUUID, destinationUUID)
		reqResult = sendReq(t, "GET", url, nil, auth)
		if balance := reqResult["balance"].(float64); balance != 70 {
			t.Fatalf("wrong balance :%v, exp:%v", balance, 70)
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
