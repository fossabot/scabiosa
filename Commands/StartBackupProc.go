package Commands

import (
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"os"
	"scabiosa/Compressor"
	"scabiosa/Logging"
	"scabiosa/SQL"
	"scabiosa/StorageTypes"
	"scabiosa/Tools"
	"time"
)

func NewStartBackupProc() *cli.Command {
	logger := Logging.Logger("backup")

	return &cli.Command{
		Name:        "backup",
		Usage:       "Starts backup process",
		Description: "Compresses and uploads/stores the backups",
		HelpName:    "backup",
		Action: func(c *cli.Context) error {
			Tools.CheckIfConfigExists()
			config := Tools.GetConfig()

			SQL.CreateDefaultTables(SQL.GetSQLInstance())

			for _, backupItem := range config.FolderToBackup {

				var storage StorageTypes.Storage
				var destPath string

				if backupItem.RemoteStorageType != "none" {
					storage = StorageTypes.CheckStorageType(backupItem.RemoteStorageType)
					destPath = checkTmpPath(backupItem.CreateLocalBackup, backupItem.LocalTargetPath)
				} else {
					destPath = backupItem.LocalTargetPath
				}

				bakFile := Compressor.CreateBakFile(backupItem.BackupName+getTimeSuffix(), backupItem.FolderPath, destPath, backupItem.BackupName)

				if backupItem.RemoteStorageType != "none" {
					StorageTypes.UploadFile(storage, bakFile, backupItem.BackupName, backupItem.RemoteTargetPath)
				}

				if !backupItem.CreateLocalBackup && backupItem.RemoteStorageType != "none" {
					backupItem.LocalTargetPath = "NONE"

					_ = os.Remove(bakFile)
					SQL.NewLogEntry(SQL.GetSQLInstance(), uuid.New(), SQL.LogInfo, backupItem.BackupName, SQL.SQLStage_DeleteTmp, SQL.REMOTE_NONE, "Deleted tmp file", time.Now())
				}

				if backupItem.RemoteStorageType == "none" {
					backupItem.CreateLocalBackup = true
					backupItem.RemoteTargetPath = "NONE"
				}
				SQL.NewBackupEntry(SQL.GetSQLInstance(), backupItem.BackupName, time.Now(), backupItem.CreateLocalBackup, backupItem.FolderPath, StorageTypes.CheckRemoteStorageType(backupItem.RemoteStorageType), backupItem.RemoteTargetPath, backupItem.LocalTargetPath)
			}

			return nil
		},
		OnUsageError: func(cc *cli.Context, err error, isSubcommand bool) error {
			if err != nil {
				logger.Fatal(err)
			}
			return err
		},
	}
}

func getTimeSuffix() string {
	currTime := time.Now()

	return "_" + currTime.Format("02-01-2006_15-04")
}

func checkTmpPath(createLocalBackup bool, targetPath string) string {
	logger := Logging.DetailedLogger("mainThread", "checkTmpPath")
	if !createLocalBackup {
		if _, err := os.Stat("tmp"); os.IsNotExist(err) {
			dirErr := os.Mkdir("tmp", 0775)
			if dirErr != nil {
				logger.Fatal(err)
			}
		}
		return "tmp"
	}

	return targetPath
}
