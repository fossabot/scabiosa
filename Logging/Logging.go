package Logging

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func checkIfLogFolderExists(){
	_, dirErr := os.Stat("logs")
	if dirErr != nil{
		log.Fatal(dirErr)
	}

	permMode, _ := os.Stat("logs")


	dirCreateErr := os.Mkdir("logs", permMode.Mode().Perm())
	if dirCreateErr != nil{
		log.Fatal(dirCreateErr)
	}

}

var logger = createLogger(logrus.WarnLevel)


func DetailedLogger(name string, section string) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"loggerName": name,
		"section": section,
	})
}

func Logger(name string) *logrus.Entry {
	return logger.WithField("loggerName", name)
}

func createLogger(logLevel logrus.Level) *logrus.Logger{
	var logger = logrus.New()

	logger.Formatter = new(logrus.TextFormatter)
	logger.Formatter.(*logrus.TextFormatter).DisableColors = false
	logger.Formatter.(*logrus.TextFormatter).DisableTimestamp = false
	logger.Level = logLevel
	logger.Out = os.Stdout

	/*file, err := os.OpenFile("logs/" + loggerName + "_" + dt.Format("02-01-2006_15_04_05") + ".log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil{
		log.Fatal(err)
	}
	logger.Out = file
	*/

	return logger
}