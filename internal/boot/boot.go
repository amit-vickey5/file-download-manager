package boot

import (
	"fmt"
	"github.com/amit/file-download-manager/pkg/config"
	"github.com/amit/file-download-manager/pkg/db"
	"github.com/amit/file-download-manager/pkg/logger"
)

var (
	Config config.Config
)

func Init() error {
	configErr := config.LoadConfig(&Config)
	if configErr != nil {
		fmt.Println("error while loading configurations")
		return configErr
	}
	fmt.Println("config loaded successfully...")

	logError := logger.InitializeLogger(Config.LoggerConfig)
	if logError != nil {
		fmt.Println("error while initializing logger")
		return logError
	}
	fmt.Println("logger initialized successfully...")

	var dbError error
	db.DbWrapper, dbError = db.InitializeDB(Config.DBConfig)
	if dbError != nil {
		fmt.Println("error while connecting to DataBase")
		return dbError
	}
	db.RepoClient = &db.Repo{db.DbWrapper.Instance()}
	fmt.Println("database connected successfully...")

	return nil
}
