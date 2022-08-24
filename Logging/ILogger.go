package Logging

import (
	"github.com/google/uuid"
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
	Hostname   string
	CurrBackup string
	CurrModule string
	CurrDest   string
	Message    string
	timeStamp  time.Time
}

func GetLoggingInstance() ILogger {

	return GetBasicLogger()
}

func InitLogger(logger ILogger) {
	logger.initLogger()
}

func NewInfoEntry(logger ILogger, entry LogEntry) {
	logger.newInfoEntry(entry)
}

func NewWarnEntry(logger ILogger, entry LogEntry) {
	logger.newWarnEntry(entry)
}

func NewErrorEntry(logger ILogger, entry LogEntry) {
	logger.newErrorEntry(entry)
}

func NewFatalEntry(logger ILogger, entry LogEntry) {
	logger.newFatalEntry(entry)
}
