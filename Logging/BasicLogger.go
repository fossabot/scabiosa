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

func (*basicLogger) Info(message string) {
	fmt.Printf("[INFO] %s\n", message)
	infoLogger.Printf(message)
}

func (*basicLogger) Warn(message string) {
	fmt.Printf("[WARN] %s\n", message)
	warningLogger.Printf(message)
}

func (*basicLogger) Error(params ...interface{}) {
	fmt.Printf("[ERROR] %s\n", params)
	errorLogger.Printf("%v", params)
}

func (*basicLogger) Fatal(params ...interface{}) {
	fmt.Printf("[FATAL] %s\n", params)
	errorLogger.Fatalf("%v", params)
}

func (*basicLogger) Init() {
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

	infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	warningLogger = log.New(file, "WARN: ", log.Ldate|log.Ltime)
	errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
}
