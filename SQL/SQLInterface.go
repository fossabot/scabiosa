package SQL

import (
	"github.com/google/uuid"
	"scabiosa/Tools"
	"time"
)

type SQLService interface {
	createDefaultTables()
	newLogEntry(uuid uuid.UUID, logType LogType, backupName string, stage SQLStage, storageType RemoteStorageType, description string, timestamp time.Time)
	newBackupEntry(backupName string, lastBackup time.Time, localBackup bool, filePath string, storageType RemoteStorageType, remotePath string, localPath string)
}

func CreateDefaultTables(sqlService SQLService) {
	sqlConfig := Tools.GetSQLConfig()
	if sqlConfig.EnableSQL {
		sqlService.createDefaultTables()
	}
}

func NewLogEntry(sqlService SQLService, uuid uuid.UUID, logType LogType, backupName string, stage SQLStage, storageType RemoteStorageType, description string, timestamp time.Time) {
	sqlConfig := Tools.GetSQLConfig()
	if sqlConfig.EnableSQL {
		sqlService.newLogEntry(uuid, logType, backupName, stage, storageType, description, timestamp)
	}
}

func NewBackupEntry(sqlService SQLService, backupName string, lastBackup time.Time, localBackup bool, filePath string, storageType RemoteStorageType, remotePath string, localPath string) {
	sqlConfig := Tools.GetSQLConfig()
	if sqlConfig.EnableSQL {
		sqlService.newBackupEntry(backupName, lastBackup, localBackup, filePath, storageType, remotePath, localPath)
	}
}

func GetSQLInstance() SQLService {
	sqlConfig := Tools.GetSQLConfig()

	if !sqlConfig.EnableSQL {
		return nil
	}

	switch sqlConfig.SqlType {
	case "mariadb":
		return GetMariaDBInstance(sqlConfig)
	case "mysql":
		return GetMariaDBInstance(sqlConfig)
	case "mssql":
		return GetMSSQLInstance(sqlConfig)
	}

	return nil
}
