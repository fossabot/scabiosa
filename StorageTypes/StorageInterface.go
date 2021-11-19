package StorageTypes

type Storage interface {
	upload(fileName string)
}

func UploadFile(storage Storage, fileName string){
	storage.upload(fileName)
}

func CheckStorageType(storageType string) Storage{

	if storageType == "azure-fileshare"{
		return GetAzureStorage()
	}

	return nil
}