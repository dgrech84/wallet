package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitlab.wedeliver.com/wedeliver/wallet/rest/api/handlers"
	"gitlab.wedeliver.com/wedeliver/wallet/utils"
	"gitlab.wedeliver.com/wedeliver/wallet/utils/config"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type walletDaoMock struct {
	DB *gorm.DB
}

func (p *walletDaoMock) GetBalance(walletId string) (decimal.Decimal, error) {
	balance, _ := decimal.NewFromString("5.5")
	return balance, nil
}

func (p *walletDaoMock) CreditWallet(walletId string, credit decimal.Decimal) (decimal.Decimal, error) {

	if walletId == "7c55f670-b644-4c92-9c54-624e5b0accc4" {
		balance, _ := decimal.NewFromString("5.7")
		return balance, nil
	}
	balance, _ := decimal.NewFromString("5.5")

	return balance, nil
}

func (p *walletDaoMock) DebitWallet(walletId string, credit decimal.Decimal) (decimal.Decimal, error) {

	if walletId == "7c55f670-b644-4c92-9c54-624e5b0accc4" {
		balance, _ := decimal.NewFromString("5.2")
		return balance, nil
	}
	balance, _ := decimal.NewFromString("5.4")
	return balance, nil
}

func TestGetBalance(t *testing.T) {
	t.Run("GetBalance - status 200", func(t *testing.T) {

		cfg, err := config.NewConfig(true)
		if err != nil {
			t.Fatalf("failed to read config : %s", err.Error())
		}

		logger := utils.NewLogger(cfg)
		if err != nil {
			t.Fatalf("failed to create logger : %s", err.Error())
		}

		walletDao := new(walletDaoMock)
		walletDao.DB = nil

		publicHandler := handlers.NewWalletAPIHandler(walletDao, logger)

		w := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/api/v1/wallets/:wallet_id/balance", nil)

		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = []gin.Param{gin.Param{Key: "wallet_id", Value: "7c55f670-b644-4c92-9c54-624e5b0accc4"}}
		ctx.Request = req

		publicHandler.GetWalletBalance(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestCreditBalance(t *testing.T) {
	t.Run("CreditBalance - status 200", func(t *testing.T) {

		cfg, err := config.NewConfig(true)
		if err != nil {
			t.Fatalf("failed to read config : %s", err.Error())
		}

		logger := utils.NewLogger(cfg)
		if err != nil {
			t.Fatalf("failed to create logger : %s", err.Error())
		}

		walletDao := new(walletDaoMock)
		walletDao.DB = nil

		publicHandler := handlers.NewWalletAPIHandler(walletDao, logger)

		w := httptest.NewRecorder()

		body := "{\"credit\": \"1.0\", \"debit\":\"1.0\"}"

		req := httptest.NewRequest("GET", "/api/v1/wallets/:wallet_id/credit", strings.NewReader(body))

		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = []gin.Param{gin.Param{Key: "wallet_id", Value: "7c55f670-b644-4c92-9c54-624e5b0accc4"}}
		ctx.Request = req

		publicHandler.CreditWallet(ctx)

		assert.Equal(t, "\"5.7\"", w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("NegativeCreditAmount - status 400", func(t *testing.T) {

		cfg, err := config.NewConfig(true)
		if err != nil {
			t.Fatalf("failed to read config : %s", err.Error())
		}

		logger := utils.NewLogger(cfg)
		if err != nil {
			t.Fatalf("failed to create logger : %s", err.Error())
		}

		walletDao := new(walletDaoMock)
		walletDao.DB = nil

		publicHandler := handlers.NewWalletAPIHandler(walletDao, logger)

		w := httptest.NewRecorder()

		body := "{\"credit\": \"-1.0\", \"debit\":\"-1.0\"}"

		req := httptest.NewRequest("GET", "/api/v1/wallets/:wallet_id/credit", strings.NewReader(body))

		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = []gin.Param{gin.Param{Key: "wallet_id", Value: "7c55f670-b644-4c92-9c54-624e5b0accc4"}}
		ctx.Request = req

		publicHandler.CreditWallet(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDebitBalance(t *testing.T) {
	t.Run("DebitBalance - status 200", func(t *testing.T) {

		cfg, err := config.NewConfig(true)
		if err != nil {
			t.Fatalf("failed to read config : %s", err.Error())
		}

		logger := utils.NewLogger(cfg)
		if err != nil {
			t.Fatalf("failed to create logger : %s", err.Error())
		}

		walletDao := new(walletDaoMock)
		walletDao.DB = nil

		publicHandler := handlers.NewWalletAPIHandler(walletDao, logger)

		w := httptest.NewRecorder()

		body := "{\"credit\": \"1.0\", \"debit\":\"1.0\"}"

		req := httptest.NewRequest("GET", "/api/v1/wallets/:wallet_id/debit", strings.NewReader(body))

		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = []gin.Param{gin.Param{Key: "wallet_id", Value: "7c55f670-b644-4c92-9c54-624e5b0accc4"}}
		ctx.Request = req

		publicHandler.DebitWallet(ctx)

		assert.Equal(t, "\"5.2\"", w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("NegativeDebitAmount - status 400", func(t *testing.T) {

		cfg, err := config.NewConfig(true)
		if err != nil {
			t.Fatalf("failed to read config : %s", err.Error())
		}

		logger := utils.NewLogger(cfg)
		if err != nil {
			t.Fatalf("failed to create logger : %s", err.Error())
		}

		walletDao := new(walletDaoMock)
		walletDao.DB = nil

		publicHandler := handlers.NewWalletAPIHandler(walletDao, logger)

		w := httptest.NewRecorder()

		body := "{\"credit\": \"-1.0\", \"debit\":\"-1.0\"}"

		req := httptest.NewRequest("GET", "/api/v1/wallets/:wallet_id/credit", strings.NewReader(body))

		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = []gin.Param{gin.Param{Key: "wallet_id", Value: "7c55f670-b644-4c92-9c54-624e5b0accc4"}}
		ctx.Request = req

		publicHandler.DebitWallet(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
