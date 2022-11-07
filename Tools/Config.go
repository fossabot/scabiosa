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
	UseHashType    string `json:"useHashType"`
	FolderToBackup []struct {
		BackupName   string `json:"backupName"`
		FolderPath   string `json:"folderPath"`
		Destinations []struct {
			DestType string `json:"destType"`
			DestPath string `json:"destPath"`
		} `json:"destinations"`
	} `json:"foldersToBackup"`
}

func readConfig() []byte {
	logger := Logging.GetLoggingInstance()

	file, err := os.ReadFile("config/config.json")
	if err != nil {
		Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error(), CurrModule: "Config"})
	}
	return file
}

func readSQLConfig() []byte {
	logger := Logging.GetLoggingInstance()

	file, err := os.ReadFile("config/sql-config.json")
	if err != nil {
		Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error(), CurrModule: "Config"})
	}
	return file
}

func CheckIfConfigExists() {
	logger := Logging.GetLoggingInstance()

	if _, err := os.Stat("config/config.json"); os.IsNotExist(err) {
		_, fileErr := os.OpenFile("config/config.json", os.O_CREATE, 0600)
		if fileErr != nil {
			Logging.NewFatalEntry(logger, Logging.LogEntry{Message: fileErr.Error(), CurrModule: "Config"})
		}
		fmt.Printf("No configs detected. Please use 'scabiosa generate-config'\n")
		os.Exit(0)
	}
}

func GenerateBaseConfig() {
	logger := Logging.GetLoggingInstance()
	var baseConfig Config

	conf, err := json.MarshalIndent(baseConfig, "", "\t")
	if err != nil {
		Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error(), CurrModule: "Config"})
	}
	for _, s := range baseConfig.FolderToBackup {
		s.BackupName = ""
	}
	err = os.WriteFile("config/config.json", conf, 0600)
	if err != nil {
		Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error(), CurrModule: "Config"})
	}
}

func GenerateAzureConfig(azure AzureConfig) {
	logger := Logging.GetLoggingInstance()

	conf, err := json.MarshalIndent(azure, "", "\t")
	if err != nil {
		Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error(), CurrModule: "Config"})
	}

	err = os.WriteFile("config/azure.json", conf, 0600)
	if err != nil {
		Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error(), CurrModule: "Config"})
	}
}

func GenerateSQLConfig(sqlConfig *SQLConfig) {
	logger := Logging.GetLoggingInstance()

	conf, err := json.MarshalIndent(sqlConfig, "", "\t")
	if err != nil {
		Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error(), CurrModule: "Config"})
	}

	err = os.WriteFile("config/sql-config.json", conf, 0600)
	if err != nil {
		Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error(), CurrModule: "Config"})
	}
}

func GetSQLConfig() SQLConfig {
	logger := Logging.GetLoggingInstance()
	var sqlConfig SQLConfig

	err := json.Unmarshal(readSQLConfig(), &sqlConfig)
	if err != nil {
		Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error(), CurrModule: "Config"})
	}

	return sqlConfig
}

func GetConfig() Config {
	logger := Logging.GetLoggingInstance()
	var config Config

	err := json.Unmarshal(readConfig(), &config)
	if err != nil {
		Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error(), CurrModule: "Config"})
	}

	return config
}
