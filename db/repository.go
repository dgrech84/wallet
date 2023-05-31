package db

import (
	"gorm.io/driver/mysql"

	"github.com/sirupsen/logrus"
	"gitlab.wedeliver.com/wedeliver/wallet/db/entity"
	"gitlab.wedeliver.com/wedeliver/wallet/utils/config"
	"gorm.io/gorm"
)

type WalletRepo struct {
	DB     *gorm.DB
	logger *logrus.Logger
}

func (c *WalletRepo) dbSetup() error {

	err := c.DB.AutoMigrate(&entity.WalletEntity{})
	if err != nil {
		return err
	}

	return nil
}

func NewWalletRepo(cfg *config.WalletConfig, logger *logrus.Logger) (WalletRepo, error) {
	repo := new(WalletRepo)

	connStr := cfg.MySqlUser + ":" + cfg.MySqlPassword + "@tcp" + "(" + cfg.MySqlHost + ":" + cfg.MySqlPort + ")/" + cfg.MySqlDatabase + "?" + "parseTime=true&loc=Local"

	var gormConfig *gorm.Config
	gormConfig = &gorm.Config{}

	db, err := gorm.Open(mysql.Open(connStr), gormConfig)
	if err != nil {
		return WalletRepo{}, err
	}

	repo.DB = db
	logger.Infof("connected to database : %s", cfg.MySqlDatabase)

	if cfg.DBSetup {
		logger.Info("migration starts...")

		err = repo.dbSetup()
		if err != nil {
			logger.Errorf("failed to auto migrate db: %s", err.Error())
		}
	}

	return WalletRepo{DB: db, logger: logger}, nil
}
