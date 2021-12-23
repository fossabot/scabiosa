package Tools

import (
	"encoding/json"
	"os"
	"scabiosa/Logging"
)

type Config struct {
	LocalBackupPath string `json:"localBackupPath"`
	SQLConfig struct{
		EnableSQL bool   `json:"enableSQL"`
		SqlType   string `json:"sqlType"`
		SqlAddress string `json:"sql-address"`
		SqlPort uint16 `json:"sql-port"`
		Database string `json:"database"`
		DbUser string `json:"db-user"`
		DbPassword string `json:"db-password"`
	} `json:"sqlConfig"`
	FolderToBackup []struct{
		BackupName string `json:"backupName"`
		FolderPath string `json:"folderPath"`
		RemoteStorageType string `json:"remoteStorageType"`
		RemoteTargetPath string `json:"remoteTargetPath"`
		CreateLocalBackup bool   `json:"createLocalBackup"`
		LocalTargetPath   string `json:"LocalTargetPath"`
	} `json:"foldersToBackup"`
}

func readConfig() []byte {
	logger := Logging.DetailedLogger("ConfigHandler", "readConfig")

	file, err := os.ReadFile("config/config.json")
	if err != nil {
		logger.Fatal(err)
	}
	return file
}

func CheckIfConfigExists(){
	logger := Logging.DetailedLogger("ConfigHandler", "CheckIfConfigExists")

	if _, err := os.Stat("config/config.json"); os.IsNotExist(err){
		_, fileErr := os.OpenFile("config/config.json", os.O_CREATE, 0775)
		if fileErr != nil{
			logger.Fatal(fileErr)
		}
		generateDefaultConfig()
	}
}

func generateDefaultConfig() {
	logger := Logging.DetailedLogger("ConfigHandler", "GenerateDefaultConfig")

	var config Config
	var conf []byte

	conf, err := json.MarshalIndent(config, "", "\t")
	//conf, err := json.Marshal(config)
	if err != nil {
		logger.Fatal(err)
	}

	err = os.WriteFile("config/config.json", conf, 0755)
	if err != nil {
		logger.Fatal(err)
	}

}

func GetConfig() Config {

	logger := Logging.DetailedLogger("ConfigHandler", "GetConfig")
	var config Config

	err := json.Unmarshal(readConfig(), &config)
	if err != nil {
		logger.Fatal(err)
	}

	return config
}