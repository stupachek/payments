package main

import (
	"pay/models"
)

func main() {
	DB := models.ConnectDataBase()
	app := New(DB)
	app.Run()
}
