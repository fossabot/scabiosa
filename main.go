package main

import (
	"github.com/urfave/cli/v2"
	"os"
	"scabiosa/Commands"
	"scabiosa/Logging"
)

func main() {
	logger := Logging.GetLoggingInstance()

	Logging.InitLogger(logger)

	app := &cli.App{
		Name:  "scabiosa",
		Usage: "Backup Util",
		Authors: []*cli.Author{
			{
				Name:  "netbenix",
				Email: "netbenix@codenoodles.de",
			},
		},
		Copyright: "(c) 2021-2022 netbenix",
		Commands: []*cli.Command{
			Commands.StartBackupProcCommand(),
			Commands.GenerateNewConfigsCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error()})
	}

	Logging.NewInfoEntry(logger, Logging.LogEntry{Message: "Finished. Exiting..."})
}