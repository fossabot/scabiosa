package Commands

import (
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
	logger := Logging.GetLoggingInstance()

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
				Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error(), CurrModule: "Backup"})
			}
			return err
		},
	}
}

func StartBackupProc() {
	Tools.CheckIfConfigExists()
	config := Tools.GetConfig()
	logger := Logging.GetLoggingInstance()

	Logging.NewInfoEntry(logger, Logging.LogEntry{Message: "Entering backup util...", CurrModule: "Backup"})

	Logging.NewInfoEntry(logger, Logging.LogEntry{Message: "Creating SQL Tables if not existing", CurrModule: "Backup"})
	SQL.CreateDefaultTables(SQL.GetSQLInstance())

	checkTmpPath()

	for _, backupItem := range config.FolderToBackup {
		Logging.NewInfoEntry(logger, Logging.LogEntry{Message: "Starting backup.", CurrBackup: backupItem.BackupName, CurrModule: "Backup"})

		bakFile := Compressor.CreateBakFile(backupItem.BackupName+getTimeSuffix(), backupItem.FolderPath, backupItem.BackupName)

		for _, backupDestination := range backupItem.Destinations {
			storage := StorageTypes.CheckStorageType(backupDestination.DestType)
			StorageTypes.UploadFile(storage, bakFile, backupItem.BackupName, backupDestination.DestPath)
			hashValue, err := Tools.CalculateHashValue(bakFile, Tools.GetHashTypeFromString(config.UseHashType))
			if err != nil {
				SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogFatal, backupItem.BackupName, SQL.SqlStageCompress, SQL.RemoteNone, "NULL", err.Error(), time.Now())
				Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error(), CurrBackup: backupItem.BackupName, CurrDest: backupDestination.DestPath, CurrModule: "Backup"})
			}

			SQL.NewBackupEntry(SQL.GetSQLInstance(), backupItem.BackupName, time.Now(), SQL.RemoteNone, backupItem.FolderPath, backupDestination.DestPath, Tools.GetHashTypeFromString(config.UseHashType), hashValue)
		}

		_ = os.Remove(bakFile)
		SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogInfo, backupItem.BackupName, SQL.SqlStageFinialzing, SQL.RemoteNone, "NULL", "Finished Backup.", time.Now())
		Logging.NewInfoEntry(logger, Logging.LogEntry{Message: "Finished backup", CurrBackup: backupItem.BackupName, CurrModule: "Backup"})
	}

}

func getTimeSuffix() string {
	currTime := time.Now()

	return "_" + currTime.Format("02-01-2006_15-04")
}

// skipcq: RVV-A0005
func checkTmpPath() {
	logger := Logging.GetLoggingInstance()
	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		dirErr := os.Mkdir("tmp", 0700)
		if dirErr != nil {
			Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error(), CurrModule: "Backup"})
		}
		Logging.NewInfoEntry(logger, Logging.LogEntry{Message: "tmp folder successfully created", CurrModule: "Backup"})
	}
}
