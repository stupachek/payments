package main

import (
	"pay/controllers"
	"pay/middleware"
	"pay/models"

	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDataBase()
	r := gin.Default()

	public := r.Group("/users")

	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)

	user := public.Group("/:user_uuid").Use(middleware.Auth())
	user.GET("/hello", controllers.Hello)

	r.Run(":8080")

}
