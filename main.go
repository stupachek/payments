package main

import (
	"pay/app"
	"pay/models"
)

func main() {
	DB := models.ConnectDataBase()
	app := app.New(DB)
	app.Run(":8080")
}
