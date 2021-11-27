package main

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"scabiosa/Compressor"
	"scabiosa/Logging"
	"scabiosa/SQL"
	"scabiosa/StorageTypes"
	"scabiosa/Tools"
	"time"
)

func main(){
	config := Tools.GetConfig()

	SQL.CreateDefaultTables(SQL.GetSQLInstance())

	for _, backupItem := range config.FolderToBackup{
		storage := StorageTypes.CheckStorageType(backupItem.StorageType)
		destPath := checkTmpPath(config, backupItem.CreateLocalBackup)

		bakFile := Compressor.CreateBakFile(backupItem.BackupName + getTimeSuffix(), backupItem.FolderPath, destPath)
		fmt.Printf(bakFile)
		StorageTypes.UploadFile(storage, bakFile)

		if !backupItem.CreateLocalBackup {
			_ = os.Remove(bakFile)
			SQL.NewLogEntry(SQL.GetSQLInstance(), uuid.New(), SQL.LogInfo, backupItem.BackupName, SQL.SQLStage_DeleteTmp, SQL.REMOTE_NONE, "Deleted tmp file" ,time.Now())
		}

	}

}

//TODO Implement SQL Backup entries
//TODO SQL Log backupName is still != config.BackupMane

func getTimeSuffix() string{
	currTime := time.Now()

	return "_" + currTime.Format("02-01-2006_15-04")
}

func checkTmpPath(config Tools.Config, createLocalBackup bool) string{
	logger := Logging.DetailedLogger("mainThread", "checkTmpPath")
	if !createLocalBackup{
		if _, err := os.Stat("tmp"); os.IsNotExist(err) {
			dirErr := os.Mkdir("tmp", 600)
			if dirErr != nil {
				logger.Fatal(err)
			}
		}
		return "tmp"
	}

	return config.LocalBackupPath
}
