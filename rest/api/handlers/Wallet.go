package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"

	"gitlab.wedeliver.com/wedeliver/wallet/rest/requests"
	"gitlab.wedeliver.com/wedeliver/wallet/utils"
)

type walletDao interface {
	GetBalance(walletId string) (decimal.Decimal, error)
	CreditWallet(walletId string, credit decimal.Decimal) (decimal.Decimal, error)
	DebitWallet(walletId string, debit decimal.Decimal) (decimal.Decimal, error)
}

type WalletAPIHandler struct {
	WalletDao walletDao
	logger    *logrus.Logger
}

const (
	timeout time.Duration = 30 * time.Second
)

func NewWalletAPIHandler(walletDao walletDao, logger *logrus.Logger) *WalletAPIHandler {
	return &WalletAPIHandler{WalletDao: walletDao, logger: logger}
}

// Get wallet balance handler
func (h *WalletAPIHandler) GetWalletBalance(c *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	walletId := c.Param("wallet_id")

	balance, err := h.WalletDao.GetBalance(walletId)
	if err != nil {
		utils.HttpError(c.Writer, h.logger, fmt.Sprintf("failed to get balance: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, balance.String())
}

// Debit wallet handler
func (h *WalletAPIHandler) DebitWallet(c *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	walletId := c.Param("wallet_id")

	var amounts requests.Amounts
	err := json.NewDecoder(c.Request.Body).Decode(&amounts)
	if err != nil {
		utils.HttpError(c.Writer, h.logger, fmt.Sprintf("failed to decode request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	amountDebited, _ := decimal.NewFromString(amounts.Debit)

	if amountDebited.LessThan(decimal.NewFromFloat(0.0)) {
		utils.HttpError(c.Writer, h.logger, "failed to debit wallet: Debited amount cannot be negative", http.StatusBadRequest)
		return
	}

	balance, err := h.WalletDao.DebitWallet(walletId, amountDebited)
	if err != nil {
		utils.HttpError(c.Writer, h.logger, fmt.Sprintf("failed to debit wallet: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, balance.String())
}

// Credit wallet handler
func (h *WalletAPIHandler) CreditWallet(c *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	walletId := c.Param("wallet_id")

	var amounts requests.Amounts
	err := json.NewDecoder(c.Request.Body).Decode(&amounts)
	if err != nil {
		utils.HttpError(c.Writer, h.logger, fmt.Sprintf("failed to decode request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	amountCredited, _ := decimal.NewFromString(amounts.Credit)

	if amountCredited.LessThan(decimal.NewFromFloat(0.0)) {
		utils.HttpError(c.Writer, h.logger, "failed to credit wallet: Credited amount cannot be negative", http.StatusBadRequest)
		return
	}

	balance, err := h.WalletDao.CreditWallet(walletId, amountCredited)
	if err != nil {
		utils.HttpError(c.Writer, h.logger, fmt.Sprintf("failed to credit wallet: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, balance.String())
}
