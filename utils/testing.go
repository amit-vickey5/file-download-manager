package main

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func main() {
	_, currentFilePath, _, _ := runtime.Caller(0)
	fmt.Println(currentFilePath)
	basepath := filepath.Dir(filepath.Dir(currentFilePath))
	fmt.Println(basepath)
}