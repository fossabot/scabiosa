package main

import (
	"encoding/json"
	"os"
	"scabiosa/Logging"
)

type Config struct {
	LocalBackupPath string `json:"localBackupPath"`
	SQLConfig struct{
		SqlType string `json:"sqlType"`
		SqlAddress string `json:"sql-address"`
		SqlPort uint16 `json:"sql-port"`
		Database string `json:"database"`
		DbUser string `json:"db-user"`
		DbPassword string `json:"db-password"`
	} `json:"sqlConfig"`
	FolderToBackup []struct{
		BackupName string `json:"backupName"`
		FolderPath string `json:"folderPath"`
		StorageType string `json:"storageType"`
		CreateLocalBackup bool `json:"createLocalBackup"`
	} `json:"foldersToBackup"`
}
type Backup struct{
	backupName string
	folderPath string
	storageType string
	createLocalBackup bool
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