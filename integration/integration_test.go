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
		reqResult := post(t, "http://localhost:8080/users/register", input, nil)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("error register")
		}
		var userUUID string = reqResult["uuid"].(string)
		reqResult = post(t, "http://localhost:8080/users/login", input, nil)
		token, ok := reqResult["token"].(string)
		if !ok {
			t.Fatal("error login")
		}
		url := fmt.Sprintf("http://localhost:8080/users/%v/accounts/new", userUUID)
		auth := make(map[string]string)
		auth["Authorization"] = token
		for k := range reqResult {
			delete(reqResult, k)
		}
		reqResult = post(t, url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		sourceUUID := reqResult["uuid"].(string)
		reqResult = post(t, url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create account error")
		}
		destinationUUID := reqResult["uuid"].(string)
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/add-money", userUUID, sourceUUID)
		inputAddMoney := controllers.AddMoneyInput{
			Amount: "10000",
		}
		reqResult = post(t, url, inputAddMoney, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("add money error")
		}
		account := reqResult["account"].(map[string]any)
		money := account["balance"].(float64)
		if money != 10000 {
			t.Fatalf("balance %v, exp: %v", money, 10000)
		}
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/transactions/new", userUUID, sourceUUID)
		inputTr := controllers.TransactionInput{
			DestinationUUID: destinationUUID,
			Amount:          "0",
		}
		reqResult = post(t, url, inputTr, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("create transaction error")
		}
		trans := reqResult["transaction"].(map[string]any)
		transUUID := trans["uuid"].(string)
		log.Print(trans)
		log.Print(transUUID)
		url = fmt.Sprintf("http://localhost:8080/users/%v/accounts/%v/transactions/%v/send", userUUID, sourceUUID, transUUID)
		reqResult = post(t, url, nil, auth)
		if _, ok := reqResult["message"]; !ok {
			t.Fatal("send transaction error")
		}

	})
}

func post(t *testing.T, url string, inputT interface{}, headers map[string]string) map[string]any {
	input, err := json.Marshal(inputT)
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(input))
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
