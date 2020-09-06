package util

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

func GenerateUniqueId() string {
	curTime := time.Now().Unix()
	key := fmt.Sprintf("%v", curTime)
	return key[len(key)-7:]
}

func GetDownloadsDirectory() string {
	_, goUtilsFileLocation, _, _ := runtime.Caller(0)
	utilDirectory := filepath.Dir(goUtilsFileLocation)
	internaDirectory := filepath.Dir(utilDirectory)
	projectDirectory := filepath.Dir(internaDirectory)
	return path.Join(projectDirectory, "downloads")
}
