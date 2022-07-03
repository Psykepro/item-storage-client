package bootstrap

import (
	domain "github.com/Psykepro/item-storage-client/_domain"
	"github.com/Psykepro/item-storage-client/config"
	"github.com/Psykepro/item-storage-client/pkg/logging"
)

func Logger(cfg *config.Config) domain.Logger {
	logger := logging.NewLogger(cfg.Logger)
	logger.InitLogger()
	logger.Infof("LogLevel: %s, Mode: %s", cfg.Logger.Level, cfg.Logger.Mode)
	return logger
}
