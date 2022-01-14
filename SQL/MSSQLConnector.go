package SQL

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/google/uuid"
	"net/url"
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

func GetMSSQLInstance(sqlConfig Tools.SQLConfig) MSSQLConnector {
	var mssql MSSQLConnector

	mssql.Address = sqlConfig.SqlAddress
	mssql.Port = sqlConfig.SqlPort
	mssql.Database = sqlConfig.Database
	mssql.DbUser = sqlConfig.DbUser
	mssql.DbPassword = sqlConfig.DbPassword

	return mssql
}

func CheckIfEventLogTableExist(db *sql.DB, mssql MSSQLConnector) bool {
	rows, _ := db.Query("SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = 'dbo' AND TABLE_NAME = 'EventLog';")
	if !rows.Next() {
		return false
	}
	return true
}

func CreateMSSQLConnection(mssql MSSQLConnector) *sql.DB {
	logger := Logging.DetailedLogger("MS-SQL", "createConnection")

	query := url.Values{}
	query.Add("app name", "scabiosa")
	query.Add("database", "scabiosa-test")

	sqlSettings := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(mssql.DbUser, mssql.DbPassword),
		Host:     fmt.Sprintf("%s:%d", mssql.Address, mssql.Port),
		RawQuery: query.Encode(),
	}

	db, err := sql.Open("sqlserver", sqlSettings.String())
	if err != nil {
		logger.Fatal(err)
	}

	return db
}

func (mssql MSSQLConnector) createDefaultTables() {}
func (mssql MSSQLConnector) newLogEntry(uuid uuid.UUID, logType LogType, backupName string, stage SQLStage, storageType RemoteStorageType, description string, timestamp time.Time) {
}
func (mssql MSSQLConnector) newBackupEntry(backupName string, lastBackup time.Time, localBackup bool, filePath string, storageType RemoteStorageType, remotePath string, localPath string) {
}
