package dao

import (
	"fmt"

	"github.com/shopspring/decimal"
	"gitlab.wedeliver.com/wedeliver/wallet/db/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletDao struct {
	DB *gorm.DB
}

func NewWalletDao(db *gorm.DB) *WalletDao {
	walletDao := new(WalletDao)
	walletDao.DB = db
	return walletDao
}

func (dao *WalletDao) GetBalance(walletId string) (decimal.Decimal, error) {

	var wallet entity.WalletEntity
	result := dao.DB.Find(&wallet, "wallet_id = ?", walletId)
	if result.Error != nil {
		return decimal.Decimal{}, result.Error
	}

	return wallet.Balance, nil
}

func (dao *WalletDao) CreditWallet(walletId string, credit decimal.Decimal) (decimal.Decimal, error) {

	var wallet entity.WalletEntity

	result := dao.DB.Find(&wallet, "wallet_id = ?", walletId)
	if result.Error != nil {
		return decimal.Decimal{}, result.Error
	}

	runningBalance := wallet.Balance
	runningBalance = runningBalance.Sub(credit)
	if runningBalance.GreaterThanOrEqual(decimal.NewFromFloat32(0.0)) {
		result := dao.DB.Model(&wallet).Clauses(clause.Returning{}).Where("wallet_id = ?", walletId).Update("balance", runningBalance)
		if result.Error != nil {
			return decimal.Decimal{}, result.Error
		}

		return wallet.Balance, nil
	}

	return decimal.Decimal{}, fmt.Errorf("wallet balance cannot go below 0")
}

func (dao *WalletDao) DebitWallet(walletId string, debit decimal.Decimal) (decimal.Decimal, error) {

	var wallet entity.WalletEntity

	result := dao.DB.Find(&wallet, "wallet_id = ?", walletId)
	if result.Error != nil {
		return decimal.Decimal{}, result.Error
	}
	runningBalance := wallet.Balance
	runningBalance = runningBalance.Add(debit)

	result = dao.DB.Model(&wallet).Clauses(clause.Returning{}).Where("wallet_id = ?", walletId).Update("balance", runningBalance)
	if result.Error != nil {
		return decimal.Decimal{}, result.Error
	}

	return wallet.Balance, nil
}
