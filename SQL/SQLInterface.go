package SQL

import (
	"github.com/google/uuid"
	"time"
)

type SQLService interface {
	createDefaultTables()
	newLogEntry(uuid uuid.UUID, logType LogType, backupName string, stage SQLStage, storageType RemoteStorageType, description string, timestamp time.Time)
	newBackupEntry(uuid uuid.UUID, backupName string, lastBackup time.Time, localBackup bool, filePath string, storageType RemoteStorageType, remotePath string, durationToBackup time.Duration, hadErrors bool)
}

func CreateDefaultTables(sqlService SQLService){
	sqlService.createDefaultTables()
}

func NewLogEntry(sqlService SQLService, uuid uuid.UUID, logType LogType, backupName string, stage SQLStage, storageType RemoteStorageType, description string, timestamp time.Time){
	sqlService.newLogEntry(uuid, logType, backupName, stage, storageType, description, timestamp)
}

func NewBackupEntry(sqlService SQLService, uuid uuid.UUID, backupName string, lastBackup time.Time, localBackup bool, filePath string, storageType RemoteStorageType, remotePath string, durationToBackup time.Duration, hadErrors bool){
	sqlService.newBackupEntry(uuid, backupName, lastBackup, localBackup, filePath, storageType, remotePath, durationToBackup, hadErrors)
}