package Commands

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"scabiosa/Logging"
	"scabiosa/Tools"
)

func GenerateNewConfigsCommand() *cli.Command {
	logger := Logging.Logger("generate-configs")

	return &cli.Command{
		Name:        "generate-config",
		Usage:       "Generates the default configs",
		Description: "Creates the specified configs",
		HelpName:    "generate-config",
		Action: func(c *cli.Context) error {
			err := os.RemoveAll("config/")
			os.Mkdir("config", 0600)

			if err != nil {
				return err
			}
			var sqlConfig Tools.SQLConfig
			var input string
			var inputInt uint8
			fmt.Printf("Want to use SQL? [Y]/[N]: ")
			fmt.Scanf("%s", &input)
			if input[0] == 'Y' || input[0] == 'y' {
				sqlConfig.EnableSQL = true
				fmt.Printf("What SQL Type do you want to use?\n")
				fmt.Printf("[0] MariaDB\n")
				fmt.Printf("[1] MySQL\n")
				fmt.Printf("[2] MS-SQL\n")
				fmt.Printf("\nSelection: ")
				fmt.Scanf("%d", &inputInt)
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
				fmt.Scanf("%s", &sqlConfig.SqlAddress)
				fmt.Printf("\nSQL Port: ")
				fmt.Scanf("%d", &sqlConfig.SqlPort)
				fmt.Printf("\nSQL Database: ")
				fmt.Scanf("%s", &sqlConfig.Database)
				fmt.Printf("\nSQL User: ")
				fmt.Scanf("%s", &sqlConfig.DbUser)
				fmt.Printf("\nSQL Password: ")
				fmt.Scanf("%s", &sqlConfig.DbPassword)

				Tools.GenerateSQLConfig(sqlConfig)
				fmt.Printf("SQL config created!")

			} else {
				sqlConfig.EnableSQL = false
			}

			fmt.Printf("\n\nWhich storage do you want to use?\n")
			fmt.Printf("[0]\tNone\n")
			fmt.Printf("[1]\tAzure File Share\n")
			fmt.Printf("\nSelection: ")
			fmt.Scanf("%d", &inputInt)
			switch inputInt {
			case 0:
				fmt.Printf("Reminder: remoteStorageType = none\n")
				//Do (nearly) nothing! :D

			case 1:
				{
					var azure Tools.AzureConfig
					fmt.Printf("\n\nStorageAccount Name: ")
					fmt.Scanf("%s", &azure.StorageAccountName)
					fmt.Printf("\nStorageAccount Key: ")
					fmt.Scanf("%s", &azure.StorageAccountKey)
					fmt.Printf("\nFileshare Name: ")
					fmt.Scanf("%s", &azure.FileshareName)
					Tools.GenerateAzureConfig(azure)
					fmt.Printf("\nAzure config created!\n")
					fmt.Printf("Reminder: remoteStorageType = azure-file\n")
				}
			default:
				fmt.Printf("Invalid input!")
				os.Exit(1)
			}

			Tools.GenerateBaseConfig()
			fmt.Printf("All configs generated!\n")
			fmt.Printf("Please modify the config.json with your backup entries.\n")

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
