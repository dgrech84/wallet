package utils

import (
	"os"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
	"gitlab.wedeliver.com/wedeliver/wallet/utils/config"
)

func NewLogger(cfg *config.WalletConfig) *logrus.Logger {

	Log := logrus.New()
	Log.SetReportCaller(true)
	Log.SetOutput(os.Stdout)

	if cfg.LogStructured {
		Log.SetFormatter(logrustash.DefaultFormatter(logrus.Fields{}))
	} else {
		Log.SetFormatter(&logrus.TextFormatter{})
	}

	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		Log.Warnf("Unknown log level configured (\"%s\"), will use \"%s\"", cfg.LogLevel, Log.Level.String())
	} else {
		Log.SetLevel(level)
	}

	return Log
}
