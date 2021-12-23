package SQL

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"scabiosa/Logging"
	"scabiosa/Tools"
	"strconv"
	"time"
)

type MariaDBConnector struct {
	Address string
	Port uint16
	Database string
	DbUser string
	DbPassword string
}

func GetMariaDBInstance(config Tools.Config) MariaDBConnector {
	var mariadb MariaDBConnector

	mariadb.Address = config.SQLConfig.SqlAddress
	mariadb.Port = config.SQLConfig.SqlPort
	mariadb.Database = config.SQLConfig.Database
	mariadb.DbUser = config.SQLConfig.DbUser
	mariadb.DbPassword = config.SQLConfig.DbPassword

	return mariadb
}


func checkIfEventLogTableExist(db *sql.DB, mariadb MariaDBConnector) bool {
	rows, _ := db.Query("SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'EventLog';", mariadb.Database)
	if !rows.Next(){ return false }
	return true
}

func checkIfBackupTableExist(db *sql.DB, mariadb MariaDBConnector) bool {
	rows, _ := db.Query("SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'Backups';", mariadb.Database)
	if !rows.Next(){ return false }
	return true
}

func checkIfBackupEntryExist(db *sql.DB, mariadb MariaDBConnector, backupName string) bool {
	rows, _ := db.Query("SELECT * FROM `" + mariadb.Database + "`.Backups WHERE BackupName = '" + backupName + "';")
	if !rows.Next(){ return false; }
	return true
}

func createMariaDBConnection(mariadb MariaDBConnector) *sql.DB{
	logger := Logging.DetailedLogger("MariaDB", "createConnection")
	db, err := sql.Open("mysql", mariadb.DbUser + ":" + mariadb.DbPassword + "@(" + mariadb.Address +  ":" +strconv.Itoa(int(mariadb.Port))+ ")/" + mariadb.Database)
	if err != nil{
		logger.Fatal(err)
	}
	return db
}

func (mariadb MariaDBConnector) createDefaultTables(){
	logger := Logging.DetailedLogger("MariaDB", "createDefaultTables")

	eventLogSQL := "create table " + mariadb.Database +".EventLog(UUID text null, LogType enum ('INFO', 'WARNING', 'ERROR', 'FATAL') null, BackupName varchar(256) null, Stage enum ('COMPRESS', 'UPLOAD', 'DELETE TMP')  null, RemoteStorage enum ('AZURE-FILE', 'AZURE-BLOB', 'NONE') null, Description text null, Timestamp datetime null);"
	backupSQL := "create table " + mariadb.Database +".Backups(UUID text null, BackupName varchar(256) null, LastBackup datetime null, LocalBackup tinyint(1) null, FilePath varchar(256) null,  RemoteStorage enum ('AZURE-FILE', 'AZURE-BLOB', 'NONE') null, RemotePath varchar(256) null);"


	db := createMariaDBConnection(mariadb)

	if !checkIfBackupTableExist(db, mariadb){
		_, err := db.Exec(backupSQL)
		if err != nil{
			logger.Fatal(err)
		}
	}

	if !checkIfEventLogTableExist(db, mariadb){
		_, err := db.Exec(eventLogSQL)
		if err != nil{
			logger.Fatal(err)
		}
	}

	_ = db.Close()
}

func (mariadb MariaDBConnector) newLogEntry(uuid uuid.UUID, logType LogType, backupName string, stage SQLStage, storageType RemoteStorageType, description string, timestamp time.Time){
	logger := Logging.DetailedLogger("MariaDB", "newLogEntry")
	db := createMariaDBConnection(mariadb)

	_, err := db.Query("INSERT INTO `" + mariadb.Database + "`.EventLog VALUES (?, ?, ?, ?, ?, ?, ?);", uuid.String(), strconv.FormatInt(int64(logType), 10), backupName, stage, strconv.FormatInt(int64(storageType), 10), description ,timestamp)
	if err != nil{
		logger.Fatal(err)
	}

}


func (mariadb MariaDBConnector) newBackupEntry(backupName string, lastBackup time.Time, localBackup bool, filePath string, storageType RemoteStorageType, remotePath string){
	logger := Logging.DetailedLogger("MariaDB", "newBackupEntry")

	db := createMariaDBConnection(mariadb)

	if checkIfBackupEntryExist(db, mariadb, backupName){
		_, err := db.Query("UPDATE `" + mariadb.Database + "`.Backups SET LastBackup = ? WHERE BackupName = ?;", lastBackup, backupName)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		_, err := db.Query("INSERT INTO `" + mariadb.Database + "`.Backups VALUES (?, ?, ?, ?, ?, ?, ?);", uuid.New(), backupName, lastBackup, localBackup, filePath, strconv.FormatInt(int64(storageType), 10), remotePath)
		if err != nil {
			logger.Fatal(err)
		}
	}
}