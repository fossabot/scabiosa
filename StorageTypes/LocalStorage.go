package StorageTypes

import (
	"io"
	"os"
	"path/filepath"
	"scabiosa/Logging"
	"scabiosa/SQL"
	"time"
)

type LocalStorage struct{}

func (LocalStorage) upload(fileName, backupName, destinationPath string) {
	logger := Logging.GetLoggingInstance()

	srcFile, srcErr := os.Open(fileName)
	if srcErr != nil {
		Logging.NewFatalEntry(logger, Logging.LogEntry{Message: srcErr.Error(), CurrBackup: backupName, CurrDest: destinationPath, CurrModule: "LocalStorage"})
	}
	defer srcFile.Close()

	destFile, destErr := os.OpenFile(destinationPath+string(os.PathSeparator)+filepath.Base(fileName), os.O_CREATE|os.O_RDWR, os.FileMode(0600))
	if destErr != nil {
		SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogFatal, backupName, SQL.SqlStageUpload, SQL.RemoteNone, destinationPath, destErr.Error(), time.Now())
		Logging.NewFatalEntry(logger, Logging.LogEntry{Message: destErr.Error(), CurrBackup: backupName, CurrDest: destinationPath, CurrModule: "LocalStorager"})
	}

	defer destFile.Close()
	SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogInfo, backupName, SQL.SqlStageUpload, SQL.RemoteNone, destinationPath, "Starting copy process", time.Now())
	Logging.NewInfoEntry(logger, Logging.LogEntry{Message: "Starting copy to destination.", CurrBackup: backupName, CurrDest: destinationPath, CurrModule: "LocalStorage"})
	if _, err := io.Copy(destFile, srcFile); err != nil {
		SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogFatal, backupName, SQL.SqlStageUpload, SQL.RemoteNone, destinationPath, err.Error(), time.Now())
		Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error(), CurrBackup: backupName, CurrDest: destinationPath, CurrModule: "LocalStorage"})
	}
	SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogInfo, backupName, SQL.SqlStageUpload, SQL.RemoteNone, destinationPath, "Finished copy process.", time.Now())
	Logging.NewInfoEntry(logger, Logging.LogEntry{Message: "Copy finished.", CurrBackup: backupName, CurrDest: destinationPath, CurrModule: "LocalStorage"})
}

func GetLocalStorage() LocalStorage {
	var locStorage LocalStorage

	return locStorage
}
