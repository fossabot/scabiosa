package StorageTypes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-storage-file-go/azfile"
	"github.com/google/uuid"
	"net/url"
	"os"
	"path/filepath"
	"scabiosa/Logging"
	"scabiosa/SQL"
	"strings"
	"time"
)

type AzureFileStorage struct{
	FileshareName string `json:"fileshareName"`
	TargetDirectory string `json:"targetDirectory"`
	StorageAccountName string `json:"storageAccountName"`
	StorageAccountKey string `json:"storageAccountKey"`
}


func (azure AzureFileStorage) upload(fileName string){
	logger := Logging.DetailedLogger("AzureFileStorage", "upload")

	file, err := os.Open(fileName)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()

	fileSize, err := file.Stat()
	if err != nil {
		logger.Fatal(err)
	}

	credential, err := azfile.NewSharedKeyCredential(azure.StorageAccountName, azure.StorageAccountKey)
	if err != nil{
		logger.Fatal(err)
	}

	u, _ := url.Parse(fmt.Sprintf("https://%s.file.core.windows.net/%s/%s/%s", azure.StorageAccountName, azure.FileshareName ,azure.TargetDirectory, filepath.Base(fileName)))

	fileURL := azfile.NewFileURL(*u, azfile.NewPipeline(credential, azfile.PipelineOptions{}))

	ctx := context.Background()

	fmt.Printf("[%s] Starting upload to Azure File Share...\n", strings.Trim(filepath.Base(fileName), ".bak"))
	//TODO Remove Hardcoded SQL Instance
	SQL.NewLogEntry(SQL.GetMariaDBInstance(), uuid.New(), SQL.LogInfo, filepath.Base(fileName), SQL.SQLStage_Upload, SQL.REMOTE_AZURE_FILE, "Starting upload.", time.Now())

	err = azfile.UploadFileToAzureFile(ctx, file, fileURL,
		azfile.UploadToAzureFileOptions{
		Parallelism: 3,
		FileHTTPHeaders: azfile.FileHTTPHeaders{
			CacheControl: "no-transform",
		},
		Progress: func(bytesTransferred int64){
			fmt.Printf("[%s] Uploaded %d of %d bytes.\n", strings.Trim(filepath.Base(fileName), ".bak") ,bytesTransferred, fileSize.Size())
		}})

	fmt.Printf("[%s] Upload finished.\n", strings.Trim(filepath.Base(fileName), ".bak"))
	//TODO Remove Hardcoded SQL Instance
	SQL.NewLogEntry(SQL.GetMariaDBInstance(), uuid.New(), SQL.LogInfo, filepath.Base(fileName), SQL.SQLStage_Upload, SQL.REMOTE_AZURE_FILE, "Finished upload.", time.Now())
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