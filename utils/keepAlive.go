package utils

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func KeepAliveWithSignals(logger *logrus.Logger, serversShutdown func()) {
	quit := make(chan os.Signal, 1)

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var osSignalMsg string

	switch <-quit {
	case syscall.SIGINT:
		osSignalMsg = "signal interrupt triggered."

	case syscall.SIGTERM:
		osSignalMsg = "signal termination triggered"
	}

	logger.Debug("service exited due to os signal", "os_signal", osSignalMsg)
	logger.Info("shutting down servers...")
	serversShutdown()
	logger.Info("servers stopped")
}
