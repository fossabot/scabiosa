package Compressor

import (
	"archive/tar"
	"compress/flate"
	"compress/gzip"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
	"scabiosa/Logging"
	"scabiosa/SQL"
	"time"
)

func CreateBakFile(fileName, folderPath, destinationPath, backupName string) string {
	logger := Logging.BasicLog

	pathToFile := destinationPath + string(os.PathSeparator) + fileName + ".bak"

	fileToWrite, err := os.OpenFile(pathToFile, os.O_CREATE|os.O_RDWR, os.FileMode(0600))
	if err != nil {
		logger.Fatal(err)
	}
	compress(fileToWrite, folderPath, backupName)

	SQL.NewLogEntry(SQL.GetSQLInstance(), uuid.New(), SQL.LogInfo, backupName, SQL.SQLStage_Compress, SQL.REMOTE_NONE, "File successfully written.", time.Now())

	return pathToFile
}

func compress(fileToWrite *os.File, folderPath, backupName string) {
	logger := Logging.BasicLog

	zr, _ := gzip.NewWriterLevel(fileToWrite, flate.BestCompression)
	tw := tar.NewWriter(zr)

	SQL.NewLogEntry(SQL.GetSQLInstance(), uuid.New(), SQL.LogInfo, backupName, SQL.SQLStage_Compress, SQL.REMOTE_NONE, "Start compression", time.Now())
	// skipcq: SCC-SA4009
	filepath.Walk(folderPath, func(file string, fi os.FileInfo, err error) error {
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			logger.Fatal(err)
		}

		relPath, _ := filepath.Rel(filepath.Dir(folderPath), file)

		header.Name = relPath
		if err := tw.WriteHeader(header); err != nil {
			logger.Fatal(err)
		}

		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				logger.Fatal(err)
			}

			if _, err := io.Copy(tw, data); err != nil {
				logger.Fatal(err)
			}
		}

		return nil
	})

	if err := tw.Close(); err != nil {
		logger.Fatal(err)
	}

	if err := zr.Close(); err != nil {
		logger.Fatal(err)
	}

	SQL.NewLogEntry(SQL.GetSQLInstance(), uuid.New(), SQL.LogInfo, backupName, SQL.SQLStage_Compress, SQL.REMOTE_NONE, "Compression complete.", time.Now())
}
