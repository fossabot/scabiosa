package StorageTypes

import "fmt"

type Storage interface {
	upload() error
}

func UploadFile(storage Storage){
	err := storage.upload()
	if err != nil{
		fmt.Print(err)
	}
}