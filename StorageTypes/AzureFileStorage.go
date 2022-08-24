package StorageTypes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-storage-file-go/azfile"
	"github.com/cheggaaa/pb/v3"
	"net/url"
	"os"
	"path/filepath"
	"scabiosa/Logging"
	"scabiosa/SQL"
	"scabiosa/Tools"
	"strings"
	"time"
)

type AzureFileStorage struct {
	FileshareName      string
	StorageAccountName string
	StorageAccountKey  string
}

func (azure AzureFileStorage) upload(fileName, backupName, destinationPath string) {
	logger := Logging.BasicLog

	file, err := os.Open(fileName)
	checkIfAzureError(backupName, destinationPath, err)
	defer file.Close()

	fileSize, err := file.Stat()
	checkIfAzureError(backupName, destinationPath, err)

	credential, err := azfile.NewSharedKeyCredential(azure.StorageAccountName, azure.StorageAccountKey)
	checkIfAzureError(backupName, destinationPath, err)

	u, err := url.Parse(fmt.Sprintf("https://%s.file.core.windows.net/%s/%s/%s", azure.StorageAccountName, azure.FileshareName, destinationPath, filepath.Base(fileName)))
	checkIfAzureError(backupName, destinationPath, err)

	fileURL := azfile.NewFileURL(*u, azfile.NewPipeline(credential, azfile.PipelineOptions{}))

	ctx := context.Background()

	logger.Info(fmt.Sprintf("[%s] Starting upload to Azure File Share...\n", backupName))
	SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogInfo, backupName, SQL.SqlStageUpload, SQL.RemoteAzureFile, destinationPath, "Starting upload.", time.Now())

	progressBar := pb.StartNew(int(fileSize.Size()))
	progressBar.Set(pb.Bytes, true)
	err = azfile.UploadFileToAzureFile(ctx, file, fileURL,
		azfile.UploadToAzureFileOptions{
			Parallelism: 3,
			FileHTTPHeaders: azfile.FileHTTPHeaders{
				CacheControl: "no-transform",
			},
			Progress: func(bytesTransferred int64) {
				progressBar.SetCurrent(bytesTransferred)
			}})

	if err != nil {
		SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogFatal, backupName, SQL.SqlStageUpload, SQL.RemoteAzureFile, destinationPath, err.Error(), time.Now())
		logger.Fatal(err)
	}
	progressBar.Finish()
	logger.Info(fmt.Sprintf("[%s] Upload finished.\n", strings.Trim(backupName, ".bak")))
	SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogInfo, backupName, SQL.SqlStageUpload, SQL.RemoteAzureFile, destinationPath, "Finished upload.", time.Now())
}

func readConfig() []byte {
	logger := Logging.BasicLog

	file, err := os.ReadFile("config/azure.json")
	if err != nil {
		logger.Fatal(err)
	}

	return file
}

func GetAzureStorage() AzureFileStorage {
	logger := Logging.BasicLog

	var azureConfig Tools.AzureConfig
	var azureFileShare AzureFileStorage

	jsonErr := json.Unmarshal(readConfig(), &azureConfig)
	if jsonErr != nil {
		logger.Fatal(jsonErr)
	}

	azureFileShare.StorageAccountName = azureConfig.StorageAccountName
	azureFileShare.StorageAccountKey = azureConfig.StorageAccountKey
	azureFileShare.FileshareName = azureConfig.FileshareName

	return azureFileShare
}

func checkIfAzureError(backupName, destinationPath string, err error) {
	logger := Logging.BasicLog
	if err != nil {
		SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogFatal, backupName, SQL.SqlStageUpload, SQL.RemoteAzureFile, destinationPath, err.Error(), time.Now())
		logger.Fatal(err)
	}
}
