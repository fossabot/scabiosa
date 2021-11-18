package main

import "scabiosa/StorageTypes"

func main(){
	azure := StorageTypes.GetAzureStorage()

	StorageTypes.UploadFile(azure)
}
