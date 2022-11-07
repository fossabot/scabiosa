package Logging

import (
	"github.com/google/uuid"
	"os"
	"time"
)

type ILogger interface {
	initLogger()
	newInfoEntry(entry LogEntry)
	newWarnEntry(entry LogEntry)
	newErrorEntry(entry LogEntry)
	newFatalEntry(entry LogEntry)
}

type LogEntry struct {
	uuid       uuid.UUID
	logType    LogType
	hostname   string
	CurrBackup string
	CurrModule string
	CurrDest   string
	Message    string
	timeStamp  time.Time
}

func GetLoggingInstance() ILogger {

	return getBasicLogger()
}

func InitLogger(logger ILogger) {
	logger.initLogger()
}

func NewInfoEntry(logger ILogger, entry LogEntry) {
	fillBaseEntryData(&entry)
	logger.newInfoEntry(entry)
}

func NewWarnEntry(logger ILogger, entry LogEntry) {
	fillBaseEntryData(&entry)
	logger.newWarnEntry(entry)
}

func NewErrorEntry(logger ILogger, entry LogEntry) {
	fillBaseEntryData(&entry)
	logger.newErrorEntry(entry)
}

func NewFatalEntry(logger ILogger, entry LogEntry) {
	fillBaseEntryData(&entry)
	logger.newFatalEntry(entry)
}

func fillBaseEntryData(entry *LogEntry) {
	entry.uuid, _ = uuid.NewUUID()
	entry.timeStamp = time.Now()
	entry.hostname, _ = os.Hostname()
}
