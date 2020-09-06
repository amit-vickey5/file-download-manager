package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

func LoadConfig(config *Config) error {
	var configPath string
	workDir := os.Getenv("WORKDIR")
	if workDir != "" {
		configPath = path.Join(workDir, "./config/config.toml")
	} else {
		_, thisFile, _, _ := runtime.Caller(1)
		configPath = path.Join(path.Dir(thisFile), "../../config/config.toml")
	}
	fmt.Println("Config File Path ::", configPath)
	loadErr := LoadConfigFromToml(configPath, config)
	if loadErr != nil {
		fmt.Println("Error while loading config from TOML file")
		return loadErr
	}
	return nil
}

func LoadConfigFromToml(filePath string, configuration *Config) error {
	content, readErr := ioutil.ReadFile(filePath)
	if readErr != nil {
		//log error
		return readErr
	}
	if _, tomlErr := toml.Decode(string(content), &configuration); tomlErr != nil {
		fmt.Println("Error while loading config :: ", tomlErr)
		return tomlErr
	}
	return nil
}