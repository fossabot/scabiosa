package Commands

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"scabiosa/Logging"
	"scabiosa/Tools"
)

func GenerateNewConfigsCommand() *cli.Command {
	logger := Logging.GetLoggingInstance()

	return &cli.Command{
		Name:        "generate-config",
		Usage:       "Generates the default configs",
		Description: "Creates the specified configs",
		HelpName:    "generate-config",
		Action: func(c *cli.Context) error {
			Logging.NewInfoEntry(logger, Logging.LogEntry{Message: "Starting config assistant", CurrModule: "ConfigAssistant"})
			err := os.RemoveAll("config/")
			if err != nil {
				Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error(), CurrModule: "ConfigAssistant"})
			}
			Logging.NewInfoEntry(logger, Logging.LogEntry{Message: "Deleted config folder", CurrModule: "ConfigAssistant"})

			dirCreateErr := os.Mkdir("config", 0700)
			if dirCreateErr != nil {
				Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error(), CurrModule: "ConfigAssistant"})
			}
			Logging.NewInfoEntry(logger, Logging.LogEntry{Message: "Created config folder", CurrModule: "ConfigAssistant"})

			var sqlConfig Tools.SQLConfig
			var input string
			var inputInt uint8
			fmt.Printf("Want to use SQL? [Y]/[N]: ")
			fmt.Scanf("%s\n", &input)
			if input[0] == 'Y' || input[0] == 'y' {
				sqlConfig.EnableSQL = true
				fmt.Printf("What SQL Type do you want to use?\n")
				fmt.Printf("[0] MariaDB\n")
				fmt.Printf("[1] MySQL\n")
				fmt.Printf("[2] MS-SQL\n")
				fmt.Printf("\nSelection: ")
				fmt.Scanf("%d\n", &inputInt)
				switch inputInt {
				case 0:
					sqlConfig.SqlType = "mariadb"
				case 1:
					sqlConfig.SqlType = "mysql"
				case 2:
					sqlConfig.SqlType = "mssql"
				default:
					fmt.Printf("Invalid input!")
					os.Exit(1)
				}

				fmt.Printf("\n\nSQL Address: ")
				fmt.Scanf("%s\n", &sqlConfig.SqlAddress)
				fmt.Printf("\nSQL Port: ")
				fmt.Scanf("%d\n", &sqlConfig.SqlPort)
				fmt.Printf("\nSQL Database: ")
				fmt.Scanf("%s\n", &sqlConfig.Database)
				fmt.Printf("\nSQL User: ")
				fmt.Scanf("%s\n", &sqlConfig.DbUser)
				fmt.Printf("\nSQL Password: ")
				fmt.Scanf("%s\n", &sqlConfig.DbPassword)

				Tools.GenerateSQLConfig(&sqlConfig)
				Logging.NewInfoEntry(logger, Logging.LogEntry{Message: "SQL config created", CurrModule: "ConfigAssistant"})

			} else {
				sqlConfig.EnableSQL = false
			}

			fmt.Printf("\n\nWhich storage do you want to use?\n")
			fmt.Printf("[0]\tLocal\n")
			fmt.Printf("[1]\tAzure File Share\n")
			fmt.Printf("\nSelection: ")
			fmt.Scanf("%d\n", &inputInt)
			switch inputInt {
			case 0:
				fmt.Printf("Reminder: destType = local\n")
				//Do (nearly) nothing! :D

			case 1:
				var azure Tools.AzureConfig
				fmt.Printf("\n\nStorageAccount Name: ")
				fmt.Scanf("%s\n", &azure.StorageAccountName)
				fmt.Printf("\nStorageAccount Key: ")
				fmt.Scanf("%s\n", &azure.StorageAccountKey)
				fmt.Printf("\nFileshare Name: ")
				fmt.Scanf("%s\n", &azure.FileshareName)
				Tools.GenerateAzureConfig(azure)
				fmt.Printf("\nAzure config created!\n")
				fmt.Printf("Reminder: remoteStorageType = azure-file\n")

			default:
				fmt.Printf("Invalid input!")
				os.Exit(1)
			}

			Tools.GenerateBaseConfig()
			Logging.NewInfoEntry(logger, Logging.LogEntry{Message: "All configs generated", CurrModule: "ConfigAssistant"})
			Logging.NewInfoEntry(logger, Logging.LogEntry{Message: "Please modify the config.json with your backup entries.", CurrModule: "ConfigAssistant"})

			return nil
		},
		OnUsageError: func(cc *cli.Context, err error, isSubcommand bool) error {
			if err != nil {
				Logging.NewFatalEntry(logger, Logging.LogEntry{Message: err.Error(), CurrModule: "ConfigAssistant"})
			}
			return err
		},
	}
}
