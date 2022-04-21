package StorageTypes

import "scabiosa/SQL"

type Storage interface {
	upload(fileName string, backupName string, destinationPath string)
}

func UploadFile(storage Storage, fileName string, backupName string, destinationPath string) {
	storage.upload(fileName, backupName, destinationPath)
}

func CheckStorageType(storageType string) Storage {

	if storageType == "azure-fileshare" {
		return GetAzureStorage()
	}
	if storageType == "local" {
		return GetLocalStorage()
	}

	return nil
}

func CheckRemoteStorageType(storageType string) SQL.RemoteStorageType {
	if storageType == "azure-fileshare" {
		return SQL.REMOTE_AZURE_FILE
	}
	if storageType == "local" {
		return SQL.REMOTE_NONE
	}

	return SQL.REMOTE_NONE
}
