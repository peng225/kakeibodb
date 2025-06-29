package cmd

import (
	"log"

	"kakeibodb/internal/mysql_client"
	"kakeibodb/internal/usecase"

	"github.com/spf13/cobra"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatal(err)
		}
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			log.Fatal(err)
		}
		credit, err := cmd.Flags().GetBool("credit")
		if err != nil {
			log.Fatal(err)
		}
		parentEventID, err := cmd.Flags().GetInt("parentEventID")
		if err != nil {
			log.Fatal(err)
		}
		if file == "" && dir == "" {
			log.Fatal("either file or dir must be specified.")
		} else if file != "" && dir != "" {
			log.Fatal("both file and dir cannot be specified.")
		}

		if credit {
			if file == "" {
				log.Fatal("file path must be specified for credit mode.")
			}
			if parentEventID < 0 {
				log.Fatalf("invalid parentEventID %d\n", parentEventID)
			}
			lceh := usecase.NewLoadCreditEventHandler(mysql_client.NewMySQLClient(dbName, dbPort, user))
			defer lceh.Close()
			lceh.LoadCreditEventFromFile(file, parentEventID)

		} else {
			leh := usecase.NewLoadEventHandler(mysql_client.NewMySQLClient(dbName, dbPort, user))
			defer leh.Close()
			if file != "" {
				leh.LoadEventFromFile(file)
			} else {
				leh.LoadEventFromDir(dir)
			}

		}
	},
}

func init() {
	eventCmd.AddCommand(loadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	loadCmd.Flags().StringP("file", "f", "", "Input file path")
	loadCmd.Flags().StringP("dir", "d", "", "Input directory path")
	loadCmd.Flags().BoolP("credit", "", false, "Load credit card event data")
	loadCmd.Flags().IntP("parentEventID", "", -1, "The parent event ID related to the credit events to be loaded")
}
