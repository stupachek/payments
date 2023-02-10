package main

import (
	"pay/app"
	"pay/controllers"
	"pay/core"
	"pay/models"
	"pay/repository"
)

func main() {
	DB := models.ConnectDataBase()
	userRepo := repository.NewGormUserRepo(DB)
	system := core.NewPaymentSystem(userRepo)
	controller := controllers.NewHttpController(system)
	app := app.New(controller)
	app.Run(":8080")
}
