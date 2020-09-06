package config

import (
	"github.com/amit/file-download-manager/pkg/db"
	"github.com/amit/file-download-manager/pkg/logger"
)

type Config struct {
	CoreConfig CoreConfig
	DBConfig db.DBConfig
	LoggerConfig logger.LoggerConfig
}

type CoreConfig struct {
	AppEnv          string
	ServiceName     string
	Hostname        string
	Port            string
	ShutdownTimeout int
	ShutdownDelay   int
}