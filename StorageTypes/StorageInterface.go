package StorageTypes

import "scabiosa/SQL"

type Storage interface {
	upload(fileName string, backupName string, destinationPath string)
}

func UploadFile(storage Storage, fileName, backupName, destinationPath string) {
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
		return SQL.RemoteAzureFile
	}
	if storageType == "local" {
		return SQL.RemoteNone
	}

	return SQL.RemoteNone
}
