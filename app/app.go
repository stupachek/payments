package app

import (
	"pay/controllers"
	"pay/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type App struct {
	DB     *gorm.DB
	Router *gin.Engine
}

func New(DB *gorm.DB) *App {
	r := gin.Default()
	public := r.Group("/users")
	c := controllers.Controller{
		DB: DB,
	}
	public.POST("/register", c.Register)
	public.POST("/login", c.Login)

	user := public.Group("/:user_uuid").Use(middleware.Auth(DB))
	user.GET("/hello", controllers.Hello)
	return &App{
		DB:     DB,
		Router: r,
	}
}

func (a *App) Run() {
	a.Router.Run(":8080")
}
