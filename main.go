package main

import (
	"github.com/urfave/cli/v2"
	"os"
	"scabiosa/Commands"
	"scabiosa/Logging"
)

func main() {
	logger := Logging.Logger("mainThread")

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
			Commands.StartBackupProc(),
			Commands.GenerateNewConfigs(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}

}
