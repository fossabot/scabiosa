package StorageTypes

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
	"scabiosa/Logging"
	"scabiosa/SQL"
	"time"
)

type LocalStorage struct{}

func (_ LocalStorage) upload(fileName, backupName, destinationPath string) {
	logger := Logging.BasicLog

	srcFile, srcErr := os.Open(fileName)
	if srcErr != nil {
		logger.Fatal(srcErr)
	}
	defer srcFile.Close()

	destFile, destErr := os.OpenFile(destinationPath+string(os.PathSeparator)+filepath.Base(fileName), os.O_CREATE|os.O_RDWR, os.FileMode(0600))
	if destErr != nil {
		logger.Fatal(destErr)
	}

	defer destFile.Close()
	SQL.NewLogEntry(SQL.GetSQLInstance(), uuid.New(), SQL.LogInfo, backupName, SQL.SQLStage_Upload, SQL.REMOTE_NONE, "Starting copy process", time.Now())
	logger.Info(fmt.Sprintf("[%s]Starting copy to %s", backupName, destinationPath))
	if _, err := io.Copy(destFile, srcFile); err != nil {
		logger.Fatal(err)
	}
	SQL.NewLogEntry(SQL.GetSQLInstance(), uuid.New(), SQL.LogInfo, backupName, SQL.SQLStage_Upload, SQL.REMOTE_NONE, "Finished copy process.", time.Now())
	logger.Info(fmt.Sprintf("[%s]Copy finished.", backupName))
}

func GetLocalStorage() LocalStorage {
	var locStorage LocalStorage

	return locStorage
}
