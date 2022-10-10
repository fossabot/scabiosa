package Compressor

import (
	"archive/tar"
	"compress/flate"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"scabiosa/Logging"
	"scabiosa/SQL"
	"time"
)

func CreateBakFile(fileName, folderPath, backupName string) string {
	logger := Logging.GetLoggingInstance()

	destinationPath := "tmp"

	pathToFile := destinationPath + string(os.PathSeparator) + fileName + ".bak"

	fileToWrite, err := os.OpenFile(pathToFile, os.O_CREATE|os.O_RDWR, os.FileMode(0600))
	if err != nil {
		SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogFatal, backupName, SQL.SqlStageCompress, SQL.RemoteNone, "NULL", err.Error(), time.Now())
		Logging.NewFatalEntry(logger, Logging.LogEntry{CurrModule: "Compression", CurrBackup: backupName, Message: err.Error()})
	}
	compress(fileToWrite, folderPath, backupName)

	SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogInfo, backupName, SQL.SqlStageCompress, SQL.RemoteNone, "NULL", "File successfully written.", time.Now())

	return pathToFile
}

func compress(fileToWrite *os.File, folderPath, backupName string) {
	logger := Logging.GetLoggingInstance()

	zr, _ := gzip.NewWriterLevel(fileToWrite, flate.BestCompression)
	tw := tar.NewWriter(zr)

	SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogInfo, backupName, SQL.SqlStageCompress, SQL.RemoteNone, "NULL", "Start compression", time.Now())
	// skipcq: SCC-SA4009
	filepath.Walk(folderPath, func(file string, fi os.FileInfo, err error) error {
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogFatal, backupName, SQL.SqlStageCompress, SQL.RemoteNone, "NULL", err.Error(), time.Now())
			Logging.NewFatalEntry(logger, Logging.LogEntry{CurrModule: "Compression", CurrBackup: backupName, Message: err.Error()})
		}

		relPath, _ := filepath.Rel(filepath.Dir(folderPath), file)

		header.Name = relPath
		if err := tw.WriteHeader(header); err != nil {
			SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogFatal, backupName, SQL.SqlStageCompress, SQL.RemoteNone, "NULL", err.Error(), time.Now())
			Logging.NewFatalEntry(logger, Logging.LogEntry{CurrModule: "Compression", CurrBackup: backupName, Message: err.Error()})
		}

		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogFatal, backupName, SQL.SqlStageCompress, SQL.RemoteNone, "NULL", err.Error(), time.Now())
				Logging.NewFatalEntry(logger, Logging.LogEntry{CurrModule: "Compression", CurrBackup: backupName, Message: err.Error()})
			}

			if _, err := io.Copy(tw, data); err != nil {
				SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogFatal, backupName, SQL.SqlStageCompress, SQL.RemoteNone, "NULL", err.Error(), time.Now())
				Logging.NewFatalEntry(logger, Logging.LogEntry{CurrModule: "Compression", CurrBackup: backupName, Message: err.Error()})
			}
		}

		return nil
	})

	if err := tw.Close(); err != nil {
		SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogFatal, backupName, SQL.SqlStageCompress, SQL.RemoteNone, "NULL", err.Error(), time.Now())
		Logging.NewFatalEntry(logger, Logging.LogEntry{CurrModule: "Compression", CurrBackup: backupName, Message: err.Error()})
	}

	if err := zr.Close(); err != nil {
		SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogFatal, backupName, SQL.SqlStageCompress, SQL.RemoteNone, "NULL", err.Error(), time.Now())
		Logging.NewFatalEntry(logger, Logging.LogEntry{CurrModule: "Compression", CurrBackup: backupName, Message: err.Error()})
	}

	SQL.NewLogEntry(SQL.GetSQLInstance(), Logging.LogInfo, backupName, SQL.SqlStageCompress, SQL.RemoteNone, "NULL", "Compression complete.", time.Now())
}
