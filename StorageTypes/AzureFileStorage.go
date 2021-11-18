package StorageTypes

import (
	"encoding/json"
	"errors"
	"os"
	"scabiosa/Logging"
)

type AzureFileStorage struct{
	azcopyPath string
	storageAccUrl string
	targetDirectory string
	SASKey string
}


func (azure AzureFileStorage) upload() error{
	//Do Stuff here
	return errors.New("lelek")
}

func readConfig() []byte {
	logger := Logging.DetailedLogger("AzureFileStorage", "readConfig")

	file, err := os.ReadFile("config/azure.json")
	if err != nil{
		logger.Fatal(err)
	}

	return file
}


func GetAzureStorage() AzureFileStorage {
	logger := Logging.DetailedLogger("AzureFileStorage", "GetAzureStorage")

	var azureStorage AzureFileStorage

	jsonErr := json.Unmarshal(readConfig(), &azureStorage)
	if jsonErr != nil{
		logger.Fatal(jsonErr)
	}

	return azureStorage
}