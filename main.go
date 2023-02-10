package main

import (
	"pay/app"
	"pay/core"
	"pay/models"
	"pay/repository"
)

func main() {
	DB := models.ConnectDataBase()
	userRepo := repository.NewGormUserRepo(DB)
	system := core.NewPaymentSystem(userRepo)
	app := app.New(DB)
	app.Run(":8080")
}
