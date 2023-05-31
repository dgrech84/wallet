package internalApi

import (
	"gitlab.wedeliver.com/wedeliver/wallet/rest"
	"gitlab.wedeliver.com/wedeliver/wallet/utils/config"

	"github.com/gin-gonic/gin"
)

func NewInternalAPIRouter(handler *InternalAPIHandler, cfg *config.WalletConfig) *gin.Engine {
	router := rest.NewRouter(cfg)

	// Routes
	router.GET("/health", handler.healthHandler)

	return router
}
