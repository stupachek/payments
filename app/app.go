package app

import (
	"pay/controllers"
	"pay/middleware"

	"github.com/gin-gonic/gin"
)

type App struct {
	controller controllers.Controller
	Router     *gin.Engine
}

func New(c controllers.Controller) *App {
	r := gin.Default()
	public := r.Group("/users")

	public.POST("/register", c.Register)
	public.POST("/login", c.Login)

	user := public.Group("/:user_uuid").Use(middleware.Auth(c))
	user.GET("/hello", controllers.Hello)
	user.POST("/accounts/new", c.NewAccount)
	user.GET("/accounts", c.GetAccounts)
	user.POST("/accounts/:account_uuid/transactions/new", c.NewTransaction)
	user.GET("/accounts/:account_uuid/transactions", c.GetTransactions)
	return &App{
		controller: c,
		Router:     r,
	}
}

func (a *App) Run(port string) {
	a.Router.Run(port)
}
