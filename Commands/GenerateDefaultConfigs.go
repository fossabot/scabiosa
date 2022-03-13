package Commands

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"scabiosa/Logging"
	"scabiosa/Tools"
)

func GenerateNewConfigsCommand() *cli.Command {
	logger := Logging.BasicLog

	return &cli.Command{
		Name:        "generate-config",
		Usage:       "Generates the default configs",
		Description: "Creates the specified configs",
		HelpName:    "generate-config",
		Action: func(c *cli.Context) error {
			logger.Info("Entering configuration assistant...")
			err := os.RemoveAll("config/")
			if err != nil {
				logger.Fatal(err)
			}
			logger.Info("Deleted config folder.")

			dirCreateErr := os.Mkdir("config", 0700)
			if dirCreateErr != nil {
				logger.Fatal(err)
			}
			logger.Info("Created config folder.")

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

				Tools.GenerateSQLConfig(sqlConfig)
				logger.Info("SQL config created!")

			} else {
				sqlConfig.EnableSQL = false
			}

			fmt.Printf("\n\nWhich storage do you want to use?\n")
			fmt.Printf("[0]\tNone\n")
			fmt.Printf("[1]\tAzure File Share\n")
			fmt.Printf("\nSelection: ")
			fmt.Scanf("%d\n", &inputInt)
			switch inputInt {
			case 0:
				fmt.Printf("Reminder: remoteStorageType = none\n")
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
			logger.Info("All configs generated!")
			logger.Info("Please modify the config.json with your backup entries.")

			return nil
		},
		OnUsageError: func(cc *cli.Context, err error, isSubcommand bool) error {
			if err != nil {
				logger.Fatal(err)
			}
			return err
		},
	}
}
