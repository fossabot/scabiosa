package SQL

import (
	"scabiosa/Logging"
	"scabiosa/Tools"
	"time"
)

type SQLService interface {
	createDefaultTables()
	newLogEntry(logType Logging.LogType, backupName string, stage SQLStage, storageType RemoteStorageType, destination, description string, timestamp time.Time)
	newBackupEntry(backupName string, lastBackup time.Time, storageType RemoteStorageType, sourcePath, destPath string, checksumType Tools.HashType, checksum string)
}

func CreateDefaultTables(sqlService SQLService) {
	sqlConfig := Tools.GetSQLConfig()
	if sqlConfig.EnableSQL {
		sqlService.createDefaultTables()
	}
}

func NewLogEntry(sqlService SQLService, logType Logging.LogType, backupName string, stage SQLStage, storageType RemoteStorageType, destination, description string, timestamp time.Time) {
	sqlConfig := Tools.GetSQLConfig()
	if sqlConfig.EnableSQL {
		sqlService.newLogEntry(logType, backupName, stage, storageType, destination, description, timestamp)
	}
}

func NewBackupEntry(sqlService SQLService, backupName string, lastBackup time.Time, storageType RemoteStorageType, sourcePath, destPath string, checksumType Tools.HashType, checksum string) {
	sqlConfig := Tools.GetSQLConfig()
	if sqlConfig.EnableSQL {
		sqlService.newBackupEntry(backupName, lastBackup, storageType, sourcePath, destPath, checksumType, checksum)
	}
}

func GetSQLInstance() SQLService {
	sqlConfig := Tools.GetSQLConfig()

	if !sqlConfig.EnableSQL {
		return nil
	}

	switch sqlConfig.SqlType {
	case "mariadb":
		return GetMariaDBInstance(&sqlConfig)
	case "mysql":
		return GetMariaDBInstance(&sqlConfig)
	case "mssql":
		return GetMSSQLInstance(&sqlConfig)
	}

	return nil
}
