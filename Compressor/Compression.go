package Compressor

import (
	"archive/tar"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
	"scabiosa/Logging"
	"scabiosa/SQL"
	"time"
)

func CreateBakFile(fileName string, folderPath string, destinationPath string) string {
	logger := Logging.DetailedLogger("Compression", "CreateBakFile")

	pathToFile := destinationPath + string(os.PathSeparator) + fileName + ".bak"

	fileToWrite, err := os.OpenFile(pathToFile, os.O_CREATE|os.O_RDWR, os.FileMode(600))
	if err != nil {
		logger.Fatal(err)
	}
	compress(fileToWrite, folderPath)

	SQL.NewLogEntry(SQL.GetSQLInstance(), uuid.New(), SQL.LogInfo, filepath.Base(folderPath), SQL.SQLStage_Compress, SQL.REMOTE_NONE, "File successfully written.", time.Now())
	fileToWrite.Close()
	return pathToFile
}

func compressFile(targetFile *os.File, file string, fi os.FileInfo, folderPath string) error {

	fileWriter, _ := gzip.NewWriterLevel(targetFile, flate.BestCompression)
	tw := tar.NewWriter(fileWriter)

	header, err := tar.FileInfoHeader(fi, file)
	if err != nil{
		return err
	}

	relPath, _ := filepath.Rel(filepath.Dir(folderPath), file)
	header.Name = relPath

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if !fi.IsDir(){
		data, err := os.Open(file)
		if err != nil {
			return err
		}

		fmt.Printf("[%s] Compressing: %s (%d bytes)\n", filepath.Base(folderPath) ,relPath, fi.Size())
		if _, err := io.Copy(tw, data); err != nil {
			return err
		}
	}

	if err := tw.Close(); err != nil {
		return err
	}

	if err := fileWriter.Close(); err != nil {
		return err
	}

	return nil
}

func compress(targetFile *os.File, folderPath string) {
	logger := Logging.DetailedLogger("Gzip", "compress")

	fmt.Printf("[%s] Start compression...\n", filepath.Base(folderPath))
	SQL.NewLogEntry(SQL.GetSQLInstance(), uuid.New(), SQL.LogInfo, filepath.Base(folderPath), SQL.SQLStage_Compress, SQL.REMOTE_NONE, "Start compression", time.Now())
	filepath.Walk(folderPath, func(file string, fi os.FileInfo, err error) error {

		//This delay is to ensure the files don't get a sudden "file aleady close" error
		time.Sleep(20 * time.Millisecond)
		go func() {
			err := compressFile(targetFile, file, fi, folderPath)
			if err != nil {
				logger.Fatal(err)
			}
		}()

		return nil
	})

	//Wait until all file writes all done
	time.Sleep(5 * time.Second)
	fmt.Printf("[%s] Compression Done.\n", filepath.Base(folderPath))
	SQL.NewLogEntry(SQL.GetSQLInstance(), uuid.New(), SQL.LogInfo, filepath.Base(folderPath), SQL.SQLStage_Compress, SQL.REMOTE_NONE, "Compression complete.", time.Now())
}