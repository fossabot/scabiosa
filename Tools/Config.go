package Tools

import (
	"encoding/json"
	"fmt"
	"os"
	"scabiosa/Logging"
)

type SQLConfig struct {
	EnableSQL  bool   `json:"enableSQL"`
	SqlType    string `json:"sqlType"`
	SqlAddress string `json:"sql-address"`
	SqlPort    uint16 `json:"sql-port"`
	Database   string `json:"database"`
	DbUser     string `json:"db-user"`
	DbPassword string `json:"db-password"`
}

type AzureConfig struct {
	FileshareName      string `json:"fileshareName"`
	StorageAccountName string `json:"storageAccountName"`
	StorageAccountKey  string `json:"storageAccountKey"`
}

type Config struct {
	FolderToBackup []struct {
		BackupName        string `json:"backupName"`
		FolderPath        string `json:"folderPath"`
		RemoteStorageType string `json:"remoteStorageType"`
		RemoteTargetPath  string `json:"remoteTargetPath"`
		CreateLocalBackup bool   `json:"createLocalBackup"`
		LocalTargetPath   string `json:"LocalTargetPath"`
	} `json:"foldersToBackup"`
}

func readConfig() []byte {
	logger := Logging.BasicLog

	file, err := os.ReadFile("config/config.json")
	if err != nil {
		logger.Fatal(err)
	}
	return file
}

func readSQLConfig() []byte {
	logger := Logging.BasicLog

	file, err := os.ReadFile("config/sql-config.json")
	if err != nil {
		logger.Fatal(err)
	}
	return file
}

func CheckIfConfigExists() {
	logger := Logging.BasicLog

	if _, err := os.Stat("config/config.json"); os.IsNotExist(err) {
		_, fileErr := os.OpenFile("config/config.json", os.O_CREATE, 0600)
		if fileErr != nil {
			logger.Fatal(fileErr)
		}
		fmt.Printf("No configs detected. Please use 'scabiosa generate-config'\n")
		os.Exit(0)
	}
}

func GenerateBaseConfig() {
	logger := Logging.BasicLog
	var baseConfig Config

	conf, err := json.MarshalIndent(baseConfig, "", "\t")
	if err != nil {
		logger.Fatal(err)
	}
	for _, s := range baseConfig.FolderToBackup {
		s.BackupName = ""
	}
	err = os.WriteFile("config/config.json", conf, 0600)
	if err != nil {
		logger.Fatal(err)
	}

}

func GenerateAzureConfig(azure AzureConfig) {
	logger := Logging.BasicLog

	conf, err := json.MarshalIndent(azure, "", "\t")
	if err != nil {
		logger.Fatal(err)
	}

	err = os.WriteFile("config/azure.json", conf, 0600)
	if err != nil {
		logger.Fatal(err)
	}
}

func GenerateSQLConfig(sqlConfig SQLConfig) {
	logger := Logging.BasicLog

	conf, err := json.MarshalIndent(sqlConfig, "", "\t")
	if err != nil {
		logger.Fatal(err)
	}

	err = os.WriteFile("config/sql-config.json", conf, 0600)
	if err != nil {
		logger.Fatal(err)
	}

}

func GetSQLConfig() SQLConfig {
	logger := Logging.BasicLog
	var sqlConfig SQLConfig

	err := json.Unmarshal(readSQLConfig(), &sqlConfig)
	if err != nil {
		logger.Fatal(err)
	}

	return sqlConfig
}

func GetConfig() Config {

	logger := Logging.BasicLog
	var config Config

	err := json.Unmarshal(readConfig(), &config)
	if err != nil {
		logger.Fatal(err)
	}

	return config
}
