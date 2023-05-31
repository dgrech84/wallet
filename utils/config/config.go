package config

import (
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

const serviceName = "research"

type WalletConfig struct {
	ServiceName     string
	ServiceBootTime time.Time
	Profile         string `envconfig:"PROFILE" default:"local"`
	ContextTimeout  int    `envconfig:"CONTEXT_TIMEOUT" default:"5" desc:"context_timeout"`
	GinMode         string `envconfig:"GIN_MODE" default:"debug"`

	// Logging
	LogLevel      string `envconfig:"LOG_LEVEL" default:"debug" yaml:"logLevel"`
	LogStructured bool   `envconfig:"LOG_STRUCTURED" default:"true" yaml:"logStructured"`

	// DB Setup
	DBSetup bool `envconfig:"DB_SETUP" default:"true" yaml:"dbSetup"`

	// Web Servers
	ServerAddress         string `envconfig:"SERVER_ADDRESS" default:":9999" desc:"localhost:9999"`
	InternalServerAddress string `envconfig:"INTERNAL_SERVER_ADDRESS" default:":9998" desc:"localhost:9998"`

	// MySql
	MySqlHost     string `envconfig:"SQL_HOST" default:"127.0.0.1" required:"true" desc:"127.0.0.1"`
	MySqlPort     string `envconfig:"SQL_PORT" default:"3306" required:"true" desc:"5432"`
	MySqlUser     string `envconfig:"SQL_USER" default:"root" required:"true" desc:"user"`
	MySqlPassword string `envconfig:"SQL_PASSWORD" default:"admin" required:"true" desc:"password"`
	MySqlDatabase string `envconfig:"SQL_DB" default:"wallet" required:"false" desc:"databse name"`
}

func initConfig(cfg *WalletConfig) error {
	cfgFileName := fmt.Sprintf("config.%s.yml", cfg.Profile)
	file, err := os.Open(cfgFileName)
	if err != nil {
		return fmt.Errorf("failed to read profile config from %s : %s", cfgFileName, err.Error())
	}

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(cfg)
	if err != nil {
		return fmt.Errorf("failed to parse the config from %s : %s", cfgFileName, err.Error())
	}

	return nil
}

func initDefaultConfig(cfg *WalletConfig) error {
	cfg.ServiceBootTime = time.Now()
	cfg.ServiceName = serviceName

	err := envconfig.Process(cfg.ServiceName, cfg)
	if err != nil {
		return fmt.Errorf("failed to read environment variables configuration : %s", err.Error())
	}

	return nil
}

func NewConfig(defaultConfig bool) (*WalletConfig, error) {
	cfg := new(WalletConfig)
	var err error

	if err = initDefaultConfig(cfg); err != nil {
		return nil, err
	}
	if !defaultConfig {
		if err = initConfig(cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
