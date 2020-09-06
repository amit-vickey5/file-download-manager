package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
)

var (
	Logger *log.Logger
)

func InitializeLogger(config LoggerConfig) error {
	var logPath string
	workDir := os.Getenv("WORKDIR")
	if workDir != "" {
		logPath = path.Join(workDir, "./logs.txt")
	} else {
		/*currentDir, err := os.Getwd()*/
		_, thisFile, _, _ := runtime.Caller(1)
		logPath = path.Join(path.Dir(thisFile), "../../logs.txt")
	}
	file, fileErr := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if fileErr != nil {
		fmt.Println("Error while opening log file :: ", fileErr)
		return fileErr
	}
	Logger = log.New(file, "", log.Ltime)
	return nil
}

func LogStatement(msg string, data interface{}) {
	if data != nil {
		Logger.Println(msg, data)
	} else {
		Logger.Println(msg)
	}
}
