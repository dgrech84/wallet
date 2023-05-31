package internalApi

import (
	"net/http"
	"time"

	"gitlab.wedeliver.com/wedeliver/wallet/utils/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type InternalAPIHandler struct {
	config *config.WalletConfig
	logger *logrus.Logger
}

func NewInternalAPIHandler(config *config.WalletConfig, logger *logrus.Logger) *InternalAPIHandler {
	return &InternalAPIHandler{config: config, logger: logger}
}

type health struct {
	Status      string `json:"status"`
	ServiceName string `json:"service_name"`
	BootTime    string `json:"boot_time"`
}

func (h *InternalAPIHandler) healthHandler(c *gin.Context) {

	healthRes := health{
		Status:      "running",
		ServiceName: h.config.ServiceName,
		BootTime:    h.config.ServiceBootTime.Format(time.RFC1123),
	}

	c.JSON(http.StatusOK, healthRes)
}
