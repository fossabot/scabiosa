package StorageTypes

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"scabiosa/Logging"
	"scabiosa/SQL"
	"time"
)

type LocalStorage struct{}

func (LocalStorage) upload(fileName, backupName, destinationPath string) {
	logger := Logging.BasicLog

	srcFile, srcErr := os.Open(fileName)
	if srcErr != nil {
		logger.Fatal(srcErr)
	}
	defer srcFile.Close()

	destFile, destErr := os.OpenFile(destinationPath+string(os.PathSeparator)+filepath.Base(fileName), os.O_CREATE|os.O_RDWR, os.FileMode(0600))
	if destErr != nil {
		SQL.NewLogEntry(SQL.GetSQLInstance(), SQL.LogFatal, backupName, SQL.SqlStageUpload, SQL.RemoteNone, destinationPath, destErr.Error(), time.Now())
		logger.Fatal(destErr)
	}

	defer destFile.Close()
	SQL.NewLogEntry(SQL.GetSQLInstance(), SQL.LogInfo, backupName, SQL.SqlStageUpload, SQL.RemoteNone, destinationPath, "Starting copy process", time.Now())
	logger.Info(fmt.Sprintf("[%s]Starting copy to %s", backupName, destinationPath))
	if _, err := io.Copy(destFile, srcFile); err != nil {
		SQL.NewLogEntry(SQL.GetSQLInstance(), SQL.LogFatal, backupName, SQL.SqlStageUpload, SQL.RemoteNone, destinationPath, err.Error(), time.Now())
		logger.Fatal(err)
	}
	SQL.NewLogEntry(SQL.GetSQLInstance(), SQL.LogInfo, backupName, SQL.SqlStageUpload, SQL.RemoteNone, destinationPath, "Finished copy process.", time.Now())
	logger.Info(fmt.Sprintf("[%s]Copy finished.", backupName))
}

func GetLocalStorage() LocalStorage {
	var locStorage LocalStorage

	return locStorage
}
