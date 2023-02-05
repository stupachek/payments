package main

import (
	"pay/controllers"
	"pay/models"

	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDataBase()
	r := gin.Default()

	public := r.Group("/users")

	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)

	r.Run(":8080")

}
