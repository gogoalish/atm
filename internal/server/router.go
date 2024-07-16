package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoalish/atm/internal/app"
	"go.uber.org/zap"
)

func NewRouter(accountCntrl *app.AccountsController, l *zap.Logger) *gin.Engine {
	router := gin.New()
	router.Use(RequestLogger(l))
	accounts := router.Group("/accounts")
	{
		accounts.POST("/", accountCntrl.Create)
		accounts.POST("/:id/deposit", accountCntrl.Deposit)
		accounts.POST("/:id/withdraw", accountCntrl.Withdraw)
		accounts.GET("/:id/balance", accountCntrl.Balance)
	}

	go accountCntrl.HandleDeposits()
	go accountCntrl.HandleBalances()
	go accountCntrl.HandleWithdrawals()
	return router
}
