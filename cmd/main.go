package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"gitlab.wedeliver.com/wedeliver/wallet/db"
	"gitlab.wedeliver.com/wedeliver/wallet/db/dao"
	"gitlab.wedeliver.com/wedeliver/wallet/rest"
	"gitlab.wedeliver.com/wedeliver/wallet/rest/api"
	"gitlab.wedeliver.com/wedeliver/wallet/rest/api/handlers"
	"gitlab.wedeliver.com/wedeliver/wallet/rest/internalApi"
	"gitlab.wedeliver.com/wedeliver/wallet/utils"
	"gitlab.wedeliver.com/wedeliver/wallet/utils/config"
)

func main() {

	cfg, err := config.NewConfig(false)
	if err != nil {
		log.Fatalf("failed to create config : %s", err.Error())
	}

	logger := utils.NewLogger(cfg)
	logger.Info("service initialising")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ContextTimeout)*time.Second)
	defer cancel()

	repo, err := db.NewWalletRepo(cfg, logger)
	if err != nil {
		logger.Fatalf("failed to connect to database %s : %s", cfg.MySqlDatabase, err.Error())
	}

	walletDao := dao.NewWalletDao(repo.DB)

	// Web router
	webRouter := rest.NewRouter(cfg)

	// Handlers
	handler := handlers.NewWalletAPIHandler(walletDao, logger)

	// Api routes
	api.NewWalletAPIRoutes(handler, webRouter)

	// Run web server
	webServer := new(rest.Server)
	go func() {
		if err := webServer.Run(webRouter, cfg.ServerAddress); err != nil && err != http.ErrServerClosed {
			logger.Infof("error occured while running web server : %s", err.Error())
		}
	}()

	// Internal Server
	internalHandler := internalApi.NewInternalAPIHandler(cfg, logger)
	internalRouter := internalApi.NewInternalAPIRouter(internalHandler, cfg)
	internalServer := new(rest.Server)
	go func() {
		if err := internalServer.Run(internalRouter, cfg.InternalServerAddress); err != nil &&
			err != http.ErrServerClosed {
			logger.Infof("error occured while running internal server : %s", err.Error())
		}
	}()

	utils.KeepAliveWithSignals(logger, func() {
		webServer.Shutdown(ctx)
		internalServer.Shutdown(ctx)
	})

	logger.Info("successfully uninitialised")
}
