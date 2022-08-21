package SQL

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"os"
	"scabiosa/Logging"
	"scabiosa/Tools"
	"strconv"
	"time"
)

type MariaDBConnector struct {
	Address    string
	Port       uint16
	Database   string
	DbUser     string
	DbPassword string
}

func GetMariaDBInstance(sqlConfig *Tools.SQLConfig) MariaDBConnector {
	var mariadb MariaDBConnector

	mariadb.Address = sqlConfig.SqlAddress
	mariadb.Port = sqlConfig.SqlPort
	mariadb.Database = sqlConfig.Database
	mariadb.DbUser = sqlConfig.DbUser
	mariadb.DbPassword = sqlConfig.DbPassword

	return mariadb
}

func (mariadb MariaDBConnector) checkIfEventLogTableExist(db *sql.DB) bool {
	logger := Logging.BasicLog
	rows, err := db.Query("SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'EventLog';", mariadb.Database)
	if err != nil {
		logger.Fatal("SQL", err)
	}
	return rows.Next()
}

func (mariadb MariaDBConnector) checkIfBackupTableExist(db *sql.DB) bool {
	logger := Logging.BasicLog
	rows, err := db.Query("SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'Backups';", mariadb.Database)
	if err != nil {
		logger.Fatal("SQL", err)
	}
	return rows.Next()
}

func (mariadb MariaDBConnector) checkIfBackupEntryExist(db *sql.DB, backupName, hostname, destPath string) bool {
	logger := Logging.BasicLog
	rows, err := db.Query("SELECT * FROM `"+mariadb.Database+"`.Backups WHERE Hostname = ? AND BackupName = ? AND DestinationPath = ?;", hostname, backupName, destPath)
	if err != nil {
		logger.Fatal("SQL", err)
	}
	return rows.Next()
}

func createMariaDBConnection(mariadb MariaDBConnector) *sql.DB {
	logger := Logging.BasicLog
	db, err := sql.Open("mysql", mariadb.DbUser+":"+mariadb.DbPassword+"@("+mariadb.Address+":"+strconv.Itoa(int(mariadb.Port))+")/"+mariadb.Database)
	if err != nil {
		logger.Fatal("SQL", err)
	}
	return db
}

func (mariadb MariaDBConnector) createDefaultTables() {
	logger := Logging.BasicLog

	eventLogSQL := "create table `" + mariadb.Database +
		"`.EventLog(UUID TEXT null, " +
		"LogType ENUM ('INFO', 'WARNING', 'ERROR', 'FATAL') null, " +
		"Hostname VARCHAR(256) null, " +
		"BackupName VARCHAR(256) null, " +
		"Stage ENUM ('COMPRESS', 'UPLOAD', 'FINALIZING')  null, " +
		"Storage ENUM ('AZURE-FILE', 'LOCAL') null, " +
		"Destination TEXT null, " +
		"Description TEXT null, " +
		"Timestamp DATETIME null);"

	backupSQL := "create table `" + mariadb.Database +
		"`.Backups(UUID TEXT null, " +
		"Hostname VARCHAR(256) null, " +
		"BackupName VARCHAR(256) null, " +
		"LastBackup DATETIME null, " +
		"Storage ENUM('AZURE-FILE', 'LOCAL') null, " +
		"SourcePath VARCHAR(512) null,  " +
		"DestinationPath VARCHAR(512) null, " +
		"ChecksumType ENUM ('SHA256', 'MD5') null, " +
		"Checksum TEXT null);"

	db := createMariaDBConnection(mariadb)

	if !mariadb.checkIfBackupTableExist(db) {
		_, err := db.Exec(backupSQL)
		if err != nil {
			logger.Fatal("SQL", err)
		}
	}

	if !mariadb.checkIfEventLogTableExist(db) {
		_, err := db.Exec(eventLogSQL)
		if err != nil {
			logger.Fatal("SQL", err)
		}
	}

	_ = db.Close()
}

func (mariadb MariaDBConnector) newLogEntry(logType LogType, backupName string, stage SQLStage, storageType RemoteStorageType, destination, description string, timestamp time.Time) {
	logger := Logging.BasicLog
	db := createMariaDBConnection(mariadb)

	hostname, _ := os.Hostname()

	_, err := db.Query("INSERT INTO `"+mariadb.Database+"`.EventLog VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);", uuid.New(), logType.String(), hostname, backupName, stage, strconv.FormatInt(int64(storageType), 10), destination, description, timestamp)
	if err != nil {
		logger.Fatal("SQL", err)
	}

}

func (mariadb MariaDBConnector) newBackupEntry(backupName string, lastBackup time.Time, storageType RemoteStorageType, sourcePath, destPath string, checksumType Tools.HashType, checksum string) {
	logger := Logging.BasicLog
	db := createMariaDBConnection(mariadb)

	hostname, _ := os.Hostname()

	if mariadb.checkIfBackupEntryExist(db, backupName, hostname, destPath) {
		_, err := db.Query("UPDATE `"+mariadb.Database+"`.Backups SET LastBackup = ?, Storage = ?, SourcePath = ?, ChecksumType = ?, Checksum = ? WHERE Hostname = ? AND BackupName = ? AND DestinationPath = ?;", lastBackup, strconv.FormatInt(int64(storageType), 10), sourcePath, strconv.FormatInt(int64(checksumType), 10), checksum, hostname, backupName, destPath)
		if err != nil {
			logger.Fatal("SQL", err)
		}
	} else {
		_, err := db.Query("INSERT INTO `"+mariadb.Database+"`.Backups VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);", uuid.New(), hostname, backupName, lastBackup, strconv.FormatInt(int64(storageType), 10), sourcePath, destPath, strconv.FormatInt(int64(checksumType), 10), checksum)
		if err != nil {
			logger.Fatal("SQL", err)
		}
	}
}
