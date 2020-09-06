package main

import "C"
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	DIR_SEP = string(filepath.Separator)
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	prettyPrintJson(currentDir)
}

func prettyPrintJson(currentDir string) {
	jsonInputFileName := "json-compressed.txt"
	jsonOutputFileName := "json-pretty-print.txt"
	parentDirectory := currentDir + DIR_SEP + "utils" + DIR_SEP + "json"
	jsonInputFilePath := parentDirectory + DIR_SEP + jsonInputFileName
	jsonOutputFilePath := parentDirectory + DIR_SEP + jsonOutputFileName
	fmt.Println("Reading input from file :: ", jsonInputFileName)
	dataBytes, err := ioutil.ReadFile(jsonInputFilePath)
	if err != nil {
		panic(err)
	}
	outputFile, err := os.Create(jsonOutputFilePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := outputFile.Close()
		if err != nil {
			panic(err)
		}
	}()
	fmt.Println("Pretty printing JSON to file :: ", jsonOutputFileName)
	var raw map[string]interface{}
	if err := json.Unmarshal(dataBytes, &raw); err != nil {
		panic(err)
	}
	outputBytes, err := json.MarshalIndent(raw, "", "    ")
	_, err = outputFile.Write(outputBytes)
	fmt.Println("DONE.!!!")
}