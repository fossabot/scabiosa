package main

import (
	"archive/tar"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"scabiosa/Logging"
)

func CreateBakFile(filename string, folderPath string, destinationPath string) string {
	logger := Logging.DetailedLogger("Compression", "CreateBakFile")

	var buf bytes.Buffer
	compress(folderPath, &buf)

	fileName := filename + ".bak"

	fileToWrite, err := os.OpenFile(destinationPath + string(os.PathSeparator) + fileName, os.O_CREATE|os.O_RDWR, os.FileMode(600))
	if err != nil {
		logger.Fatal(err)
	}

	if _, err := io.Copy(fileToWrite, &buf); err != nil {
		logger.Fatal(err)
	}

	return fileName
}


func compress(folderPath string, buf io.Writer){
	logger := Logging.DetailedLogger("Gzip", "compress")

	zr, _ := gzip.NewWriterLevel(buf, flate.BestCompression)
	tw := tar.NewWriter(zr)

	fmt.Printf("[%s] Start compression...\n", filepath.Base(folderPath))

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

		if !fi.IsDir(){
			data, err := os.Open(file)
			if err != nil {
				logger.Fatal(err)
			}

			fmt.Printf("[%s] Compressing: %s (%d bytes)\n", filepath.Base(folderPath) ,relPath, fi.Size())
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


	fmt.Printf("[%s] Compression Done.\n", filepath.Base(folderPath))
}