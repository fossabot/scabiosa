package SQL

import (
	"github.com/google/uuid"
	"scabiosa/Tools"
	"time"
)

type SQLService interface {
	createDefaultTables()
	newLogEntry(uuid uuid.UUID, logType LogType, backupName string, stage SQLStage, storageType RemoteStorageType, description string, timestamp time.Time)
	newBackupEntry(backupName string, lastBackup time.Time, localBackup bool, filePath string, storageType RemoteStorageType, remotePath string)
}

func CreateDefaultTables(sqlService SQLService){
	config := Tools.GetConfig()
	if config.SQLConfig.EnableSQL{
		sqlService.createDefaultTables()
	}
}

func NewLogEntry(sqlService SQLService, uuid uuid.UUID, logType LogType, backupName string, stage SQLStage, storageType RemoteStorageType, description string, timestamp time.Time){
	config := Tools.GetConfig()
	if config.SQLConfig.EnableSQL{
		sqlService.newLogEntry(uuid, logType, backupName, stage, storageType, description, timestamp)
	}
}

func NewBackupEntry(sqlService SQLService, backupName string, lastBackup time.Time, localBackup bool, filePath string, storageType RemoteStorageType, remotePath string){
	config := Tools.GetConfig()
	if config.SQLConfig.EnableSQL{
		sqlService.newBackupEntry(backupName, lastBackup, localBackup, filePath, storageType, remotePath)
	}
}

func GetSQLInstance() SQLService{
	config := Tools.GetConfig()

	if !config.SQLConfig.EnableSQL { return nil }

	switch config.SQLConfig.SqlType {
		case "mariadb": {return GetMariaDBInstance(config)}
	}

	return nil
}