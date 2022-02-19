package Logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

var logger = createLogger(logrus.WarnLevel)

func DetailedLogger(name string, section string) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"loggerName": name,
		"section":    section,
	})
}

func Logger(name string) *logrus.Entry {
	return logger.WithField("loggerName", name)
}

func createLogger(logLevel logrus.Level) *logrus.Logger {
	var logger = logrus.New()

	logger.Formatter = new(logrus.TextFormatter)
	logger.Formatter.(*logrus.TextFormatter).DisableColors = false
	logger.Formatter.(*logrus.TextFormatter).DisableTimestamp = false
	logger.Level = logLevel
	logger.Out = os.Stdout

	return logger
}
