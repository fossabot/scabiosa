package SQL

import (
	"database/sql"
	"fmt"
	mssqlpkg "github.com/denisenkom/go-mssqldb"
	"github.com/google/uuid"
	"net/url"
	"os"
	"scabiosa/Logging"
	"scabiosa/Tools"
	"time"
)

type MSSQLConnector struct {
	Address    string
	Port       uint16
	Database   string
	DbUser     string
	DbPassword string
}

func GetMSSQLInstance(sqlConfig *Tools.SQLConfig) MSSQLConnector {
	var mssql MSSQLConnector

	mssql.Address = sqlConfig.SqlAddress
	mssql.Port = sqlConfig.SqlPort
	mssql.Database = sqlConfig.Database
	mssql.DbUser = sqlConfig.DbUser
	mssql.DbPassword = sqlConfig.DbPassword

	return mssql
}

func (MSSQLConnector) checkIfEventLogTableExist(db *sql.DB) bool {
	logger := Logging.BasicLog
	rows, err := db.Query("SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = 'dbo' AND TABLE_NAME = 'EventLog';")
	if err != nil {
		logger.Fatal("SQL", err)
	}
	return rows.Next()
}

func (MSSQLConnector) checkIfBackupTableExist(db *sql.DB) bool {
	logger := Logging.BasicLog
	rows, err := db.Query("SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = 'dbo' AND TABLE_NAME = 'Backups';")
	if err != nil {
		logger.Fatal("SQL", err)
	}
	return rows.Next()
}

func (MSSQLConnector) checkIfBackupEntryExist(db *sql.DB, backupName, hostname, destPath string) bool {
	logger := Logging.BasicLog
	query := fmt.Sprintf("SELECT * FROM dbo.Backups WHERE Hostname = '%s' AND BackupName = '%s' AND DestinationPath = '%s'", hostname, backupName, destPath)
	rows, err := db.Query(query)
	if err != nil {
		logger.Fatal("SQL", err)
	}
	return rows.Next()
}

func createMSSQLConnection(mssql MSSQLConnector) *sql.DB {
	logger := Logging.BasicLog

	query := url.Values{}
	query.Add("app name", "scabiosa")
	query.Add("database", mssql.Database)

	sqlSettings := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(mssql.DbUser, mssql.DbPassword),
		Host:     fmt.Sprintf("%s:%d", mssql.Address, mssql.Port),
		RawQuery: query.Encode(),
	}

	connector, err := mssqlpkg.NewConnector(sqlSettings.String())
	if err != nil {
		logger.Fatal("SQL", err)
	}

	connector.SessionInitSQL = "SET ANSI_NULLS ON"

	db := sql.OpenDB(connector)

	return db
}

func (mssql MSSQLConnector) createDefaultTables() {
	logger := Logging.BasicLog

	eventLogSQL := "create table dbo.EventLog(" +
		"UUID TEXT null, " +
		"LogType ENUM ('INFO', 'WARNING', 'ERROR', 'FATAL') null, " +
		"Hostname VARCHAR(256) null, " +
		"BackupName VARCHAR(256) null, " +
		"Stage ENUM ('COMPRESS', 'UPLOAD', 'FINALIZING')  null, " +
		"Storage ENUM ('AZURE-FILE', 'LOCAL') null, " +
		"Destination TEXT null, " +
		"Description TEXT null, " +
		"Timestamp DATETIME null);"

	backupSQL := "create table dbo.Backups(" +
		"UUID TEXT null, " +
		"Hostname VARCHAR(256) null, " +
		"BackupName VARCHAR(256) null, " +
		"LastBackup DATETIME null, " +
		"Storage ENUM('AZURE-FILE', 'LOCAL') null, " +
		"SourcePath VARCHAR(512) null,  " +
		"DestinationPath VARCHAR(512) null, " +
		"ChecksumType ENUM ('SHA256', 'MD5') null, " +
		"Checksum TEXT null);"

	db := createMSSQLConnection(mssql)

	if !mssql.checkIfBackupTableExist(db) {
		_, err := db.Exec(backupSQL)
		if err != nil {
			logger.Fatal("SQL", err)
		}
	}

	if !mssql.checkIfEventLogTableExist(db) {
		_, err := db.Exec(eventLogSQL)
		if err != nil {
			logger.Fatal("SQL", err)
		}
	}

	_ = db.Close()
}

// skipcq: RVV-A0005
func (mssql MSSQLConnector) newLogEntry(logType LogType, backupName string, stage SQLStage, storageType RemoteStorageType, destination, description string, timestamp time.Time) {
	logger := Logging.BasicLog
	db := createMSSQLConnection(mssql)

	hostname, _ := os.Hostname()
	query := fmt.Sprintf("INSERT INTO dbo.EventLog VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')", uuid.New(), logType.String(), hostname, backupName, stage.String(), storageType.String(), destination, description, timestamp.Format("2006-01-02 15:04:05.999"))
	_, err := db.Query(query)
	if err != nil {
		logger.Fatal("SQL", err)
	}
}
func (mssql MSSQLConnector) newBackupEntry(backupName string, lastBackup time.Time, storageType RemoteStorageType, sourcePath, destPath string, checksumType Tools.HashType, checksum string) {
	logger := Logging.BasicLog
	db := createMSSQLConnection(mssql)

	hostname, _ := os.Hostname()

	if mssql.checkIfBackupEntryExist(db, backupName, hostname, destPath) {
		queryUpdate := fmt.Sprintf("UPDATE dbo.Backups SET Lastbackup = '%s', Storage = '%s', SourcePath = '%s', ChecksumType = '%s', Checksum = '%s' WHERE Hostname = '%s' AND BackupName = '%s' AND DestinationPath = '%s'", lastBackup.Format("2006-01-02 15:04:05.999"), storageType.String(), sourcePath, checksumType.String(), checksum, hostname, backupName, destPath)
		_, err := db.Query(queryUpdate)
		if err != nil {
			logger.Fatal("SQL", err)
		}
	} else {
		queryInsert := fmt.Sprintf("INSERT INTO dbo.Backups VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s')", uuid.New(), hostname, backupName, lastBackup.Format("2006-01-02 15:04:05.999"), storageType.String(), sourcePath, destPath, checksumType.String(), checksum)
		_, err := db.Query(queryInsert)
		if err != nil {
			logger.Fatal("SQL", err)
		}
	}
}
