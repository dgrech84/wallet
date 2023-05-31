package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.wedeliver.com/wedeliver/wallet/rest/api/handlers"
)

func NewWalletAPIRoutes(handler *handlers.WalletAPIHandler, router *gin.Engine) {

	// Fields  API routes
	fieldsGroup := router.Group("/api/v1/wallets")
	{

		fieldsGroup.GET("/:wallet_id/balance", handler.GetWalletBalance)
		fieldsGroup.POST("/:wallet_id/credit", handler.CreditWallet)
		fieldsGroup.POST("/:wallet_id/debit", handler.DebitWallet)
	}
}
