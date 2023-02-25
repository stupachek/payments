package main

import (
	"log"
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
	err := controller.System.SetupAdmin()
	if err != nil {
		log.Fatalf("can't create admin, err %v", err.Error())
	}
	app := app.New(controller)
	app.Run(":8080")
}
