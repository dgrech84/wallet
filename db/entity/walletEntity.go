package entity

import (
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type WalletEntity struct {
	WalletID uuid.UUID       `json:"wallet_id" gorm:"type:char(36); primaryKey; not null"`
	Name     string          `json:"name"`
	Surname  string          `json:"surname"`
	Balance  decimal.Decimal `json:"balance" gorm:"type:decimal(32,18)"`
}

func (WalletEntity) TableName() string {
	return "walletEntity"
}
