package Logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

type basicLogger struct{}

var (
	warningLogger *log.Logger
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	BasicLog      *basicLogger
)

func (logger *basicLogger) Info(message string) {
	fmt.Printf("[INFO] %s\n", message)
	infoLogger.Printf(message)
}

func (logger *basicLogger) Warn(message string) {
	fmt.Printf("[WARN] %s\n", message)
	warningLogger.Printf(message)
}

func (logger *basicLogger) Error(err error) {
	fmt.Printf("[ERROR] %s\n", err.Error())
	errorLogger.Print(err)
}

func (logger *basicLogger) Fatal(err error) {
	fmt.Printf("[ERROR] %s\n", err.Error())
	errorLogger.Fatal(err)
}

func (logger *basicLogger) Init() {
	_, dirErr := os.Stat("logs")
	if dirErr != nil {
		dirCreateErr := os.Mkdir("logs", 0700)
		if dirCreateErr != nil {
			log.Fatal(dirCreateErr)
		}
	}

	dt := time.Now().Local()

	file, err := os.OpenFile("logs/"+dt.Format("02-01-2006_15-04")+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}

	infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	warningLogger = log.New(file, "WARN: ", log.Ldate|log.Ltime)
	errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
}
