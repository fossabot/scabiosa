package Commands

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"scabiosa/Compressor"
	"scabiosa/Logging"
	"scabiosa/SQL"
	"scabiosa/StorageTypes"
	"scabiosa/Tools"
	"time"
)

func StartBackupProcCommand() *cli.Command {
	logger := Logging.BasicLog

	return &cli.Command{
		Name:        "backup",
		Usage:       "Starts backup process",
		Description: "Compresses and uploads/stores the backups",
		HelpName:    "backup",
		Action: func(c *cli.Context) error {
			StartBackupProc()
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

func StartBackupProc() {
	Tools.CheckIfConfigExists()
	config := Tools.GetConfig()
	logger := Logging.BasicLog

	logger.Info("Entering backup util...")

	logger.Info("Creating SQL Tables if not existing")
	SQL.CreateDefaultTables(SQL.GetSQLInstance())

	checkTmpPath()

	for _, backupItem := range config.FolderToBackup {
		logger.Info(fmt.Sprintf("Starting backup for %s", backupItem.BackupName))

		bakFile := Compressor.CreateBakFile(backupItem.BackupName+getTimeSuffix(), backupItem.FolderPath, backupItem.BackupName)

		for _, backupDestination := range backupItem.Destinations {
			storage := StorageTypes.CheckStorageType(backupDestination.DestType)
			StorageTypes.UploadFile(storage, bakFile, backupItem.BackupName, backupDestination.DestPath)
			hashValue, err := Tools.CalculateHashValue(bakFile, Tools.GetHashTypeFromString(config.UseHashType))
			if err != nil {
				SQL.NewLogEntry(SQL.GetSQLInstance(), SQL.LogFatal, backupItem.BackupName, SQL.SqlStageCompress, SQL.RemoteNone, "NULL", err.Error(), time.Now())
				logger.Fatal(err)
			}

			SQL.NewBackupEntry(SQL.GetSQLInstance(), backupItem.BackupName, time.Now(), SQL.RemoteNone, backupItem.FolderPath, backupDestination.DestPath, Tools.GetHashTypeFromString(config.UseHashType), hashValue)
		}

		_ = os.Remove(bakFile)
		SQL.NewLogEntry(SQL.GetSQLInstance(), SQL.LogInfo, backupItem.BackupName, SQL.SqlStageFinialzing, SQL.RemoteNone, "NULL", "Finished Backup.", time.Now())
		logger.Info(fmt.Sprintf("Finished backup for %s", backupItem.BackupName))
	}

}

func getTimeSuffix() string {
	currTime := time.Now()

	return "_" + currTime.Format("02-01-2006_15-04")
}

// skipcq: RVV-A0005
func checkTmpPath() {
	logger := Logging.BasicLog
	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		dirErr := os.Mkdir("tmp", 0700)
		if dirErr != nil {
			logger.Fatal(err)
		}
		logger.Info("tmp folder successfully created.")
	}
}
