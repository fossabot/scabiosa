package main

import (
	"encoding/json"
	"os"
	"scabiosa/Logging"
)

type Config struct {
	_7zPath string
	_7zArgs string
	storageType string
	localbackup uint8
	localbackupPath string
}

func readConfig() []byte {
	logger := Logging.DetailedLogger("ConfigHandler", "readConfig")

	file, err := os.ReadFile("config/config.json")
	if err != nil {
		logger.Fatal(err)
	}

	return file
}

func GetConfig() Config {
	logger := Logging.DetailedLogger("ConfigHandler", "GetConfig()")
	var config Config

	err := json.Unmarshal(readConfig(), &config)
	if err != nil {
		logger.Fatal(err)
	}

	return config
}