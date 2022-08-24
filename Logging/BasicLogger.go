package Logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

type BasicLogger struct{}

func GetBasicLogger() BasicLogger {
	var basicLogger BasicLogger

	return basicLogger
}

var (
	warnLog  *log.Logger
	infoLog  *log.Logger
	errorLog *log.Logger
)

func (_ BasicLogger) initLogger() {
	_, dirErr := os.Stat("logs")
	if dirErr != nil {
		dirCreateErr := os.Mkdir("logs", 0700)
		if dirCreateErr != nil {
			// skipcq: RVV-A0003
			log.Fatal(dirCreateErr)
		}
	}

	dt := time.Now().Local()
	file, err := os.OpenFile("logs/"+dt.Format("02-01-2006_15-04")+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		// skipcq: RVV-A0003
		log.Fatal(err)
	}

	infoLog = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	warnLog = log.New(file, "WARN: ", log.Ldate|log.Ltime)
	errorLog = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
}

func (_ BasicLogger) newInfoEntry(entry LogEntry) {
	fmt.Printf("[INFO][%s] %s\n", entry.CurrBackup, entry.Message)
	infoLog.Printf("[%s] %s", entry.CurrBackup, entry.Message)
}

func (_ BasicLogger) newWarnEntry(entry LogEntry) {
	fmt.Printf("[WARN][%s] %s\n", entry.CurrBackup, entry.Message)
	warnLog.Printf("[%s] %s", entry.CurrBackup, entry.Message)
}

func (_ BasicLogger) newErrorEntry(entry LogEntry) {
	fmt.Printf("[ERROR][%s] %s\n", entry.CurrBackup, entry.Message)
	errorLog.Printf("[%s] %s", entry.CurrBackup, entry.Message)
}

func (_ BasicLogger) newFatalEntry(entry LogEntry) {
	fmt.Printf("[FATAL][%s] %s\n", entry.CurrBackup, entry.Message)
	errorLog.Fatalf("[%s] %s", entry.CurrBackup, entry.Message)
}
