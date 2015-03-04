package controller

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	TokenKey       string
	Salt           string
	EnableRegister bool
	MongoDb        struct {
		Host     string
		Port     int
		Database string
		User     string
		Password string
	}
}

func ReadConfig(path string) (Config, error) {
	configFile, readErr := readFile(path)
	if readErr != nil {
		return Config{}, readErr
	}
	return parseConfig(configFile)
}

func readFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func parseConfig(file []byte) (Config, error) {
	var config Config
	err := json.Unmarshal(file, &config)
	return config, err
}
