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
	rows, _ := db.Query("SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'EventLog';", mariadb.Database)
	return rows.Next()
}

func (mariadb MariaDBConnector) checkIfBackupTableExist(db *sql.DB) bool {
	rows, _ := db.Query("SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'Backups';", mariadb.Database)
	return rows.Next()
}

func (mariadb MariaDBConnector) checkIfBackupEntryExist(db *sql.DB, backupName, hostname string) bool {
	rows, _ := db.Query("SELECT * FROM `"+mariadb.Database+"`.Backups WHERE Hostname = ? AND BackupName = ?;", hostname, backupName)
	return rows.Next()
}

func createMariaDBConnection(mariadb MariaDBConnector) *sql.DB {
	logger := Logging.BasicLog
	db, err := sql.Open("mysql", mariadb.DbUser+":"+mariadb.DbPassword+"@("+mariadb.Address+":"+strconv.Itoa(int(mariadb.Port))+")/"+mariadb.Database)
	if err != nil {
		logger.Fatal(err)
	}
	return db
}

func (mariadb MariaDBConnector) createDefaultTables() {
	logger := Logging.BasicLog

	eventLogSQL := "create table `" + mariadb.Database + "`.EventLog(UUID text null, LogType enum ('INFO', 'WARNING', 'ERROR', 'FATAL') null, Hostname varchar(256) null,BackupName varchar(256) null, Stage enum ('COMPRESS', 'UPLOAD', 'DELETE TMP')  null, RemoteStorage enum ('AZURE-FILE', 'AZURE-BLOB', 'NONE') null, Description text null, Timestamp datetime null);"
	backupSQL := "create table `" + mariadb.Database + "`.Backups(UUID text null, Hostname varchar(256) null, BackupName varchar(256) null, LastBackup datetime null, LocalBackup tinyint(1) null, FilePath varchar(256) null, RemoteStorage enum ('AZURE-FILE', 'AZURE-BLOB', 'NONE') null, RemotePath varchar(256) null, LocalPath varchar(256) null);"

	db := createMariaDBConnection(mariadb)

	if !mariadb.checkIfBackupTableExist(db) {
		_, err := db.Exec(backupSQL)
		if err != nil {
			logger.Fatal(err)
		}
	}

	if !mariadb.checkIfEventLogTableExist(db) {
		_, err := db.Exec(eventLogSQL)
		if err != nil {
			logger.Fatal(err)
		}
	}

	_ = db.Close()
}

func (mariadb MariaDBConnector) newLogEntry(uuid uuid.UUID, logType LogType, backupName string, stage SQLStage, storageType RemoteStorageType, description string, timestamp time.Time) {
	logger := Logging.BasicLog
	db := createMariaDBConnection(mariadb)

	hostname, _ := os.Hostname()

	_, err := db.Query("INSERT INTO `"+mariadb.Database+"`.EventLog VALUES (?, ?, ?, ?, ?, ?, ?, ?);", uuid.String(), strconv.FormatInt(int64(logType), 10), hostname, backupName, stage, strconv.FormatInt(int64(storageType), 10), description, timestamp)
	if err != nil {
		logger.Fatal(err)
	}

}

func (mariadb MariaDBConnector) newBackupEntry(backupName string, lastBackup time.Time, localBackup bool, filePath string, storageType RemoteStorageType, remotePath, localPath string) {
	logger := Logging.BasicLog
	db := createMariaDBConnection(mariadb)

	hostname, _ := os.Hostname()

	if mariadb.checkIfBackupEntryExist(db, backupName, hostname) {
		_, err := db.Query("UPDATE `"+mariadb.Database+"`.Backups SET LastBackup = ?, LocalBackup = ?, RemoteStorage = ?, RemotePath = ?, LocalPath = ? WHERE Hostname = ? AND BackupName = ?;", lastBackup, localBackup, strconv.FormatInt(int64(storageType), 10), remotePath, localPath, hostname, backupName)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		_, err := db.Query("INSERT INTO `"+mariadb.Database+"`.Backups VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);", uuid.New(), hostname, backupName, lastBackup, localBackup, filePath, strconv.FormatInt(int64(storageType), 10), remotePath, localPath)
		if err != nil {
			logger.Fatal(err)
		}
	}
}
