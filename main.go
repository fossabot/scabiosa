package main

import (
	"os"
	"scabiosa/Logging"
	"scabiosa/StorageTypes"
	"time"
)

func main(){
	config := GetConfig()

	for _, backupItem := range config.FolderToBackup{
		storage := StorageTypes.CheckStorageType(backupItem.StorageType)
		destPath := checkTmpPath(config, backupItem.CreateLocalBackup)

		bakFile := CreateBakFile(backupItem.BackupName + getTimeSuffix(), backupItem.FolderPath, destPath)

		StorageTypes.UploadFile(storage, destPath + string(os.PathSeparator) + bakFile)

		if !backupItem.CreateLocalBackup {
			_ = os.Remove(destPath + string(os.PathSeparator) + bakFile)
		}

	}

}


func getTimeSuffix() string{
	currTime := time.Now()

	return "_" + currTime.Format("02-01-2006_15-04")
}

func checkTmpPath(config Config, createLocalBackup bool) string{
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
