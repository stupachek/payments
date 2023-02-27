package app

import (
	"context"
	"net/http"
	"payment/controllers"
	"payment/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	controller controllers.Controller
	Router     *gin.Engine
	Server     http.Server
}

func New(c controllers.Controller) *App {
	r := gin.Default()
	public := r.Group("/users")
	public.POST("/register", c.Register)
	public.Use(middleware.CheckBlockedUser(c))
	public.POST("/login", c.Login)
	user := public.Group("/:user_uuid")
	user.Use(middleware.Auth(c))
	admin := r.Group("/admin/:user_uuid")
	admin.Use(middleware.Auth(c), middleware.CheckAdmin(c))
	admin.POST("/update-role", c.ChangeRole)
	admin.POST("users/:target_uuid/block", c.BlockUser)
	admin.POST("users/:target_uuid/unblock", c.UnblockUser)
	admin.POST("/accounts/:account_uuid/unblock", c.UnblockAccount)
	admin.GET("/accounts/requested", c.GetAccountsRequested)
	user.POST("/accounts/new", c.NewAccount)
	user.GET("/accounts", c.GetAccounts)
	account := user.Group("/accounts/:account_uuid")
	account.Use(middleware.CheckAccount(c))
	account.GET("", c.GetAccount)
	account.POST("/block", c.BlockAccount)
	account.POST("/unblock", c.RequestUnblockAccount)
	account.Use(middleware.CheckBlockedAccount(c))
	account.POST("/transactions/new", c.NewTransaction)
	account.GET("/transactions", c.GetTransactions)
	account.POST("/add-money", c.AddMoney)
	account.POST("/transactions/:transaction_uuid/send", c.SendTransaction)
	return &App{
		controller: c,
		Router:     r,
	}
}

func (a *App) Run(port string) error {
	server := http.Server{Addr: port, Handler: a.Router}
	a.Server = server
	return a.Server.ListenAndServe()
}

func (a *App) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.Server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
