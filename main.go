package main

import (
	"payment/app"
	"payment/controllers"
	"payment/core"
	"payment/repository"
)

func main() {
	DB := repository.ConnectDataBase()
	userRepo := repository.NewGormUserRepo(DB)
	system := core.NewPaymentSystem(userRepo)
	controller := controllers.NewHttpController(system)
	app := app.New(controller)
	app.Run(":8080")
}
