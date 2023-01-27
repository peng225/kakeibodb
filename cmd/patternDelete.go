package cmd

import (
	"kakeibodb/mysql_client"
	"kakeibodb/usecase"
	"log"

	"github.com/spf13/cobra"
)

// patternDeleteCmd represents the delete command
var patternDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		patternID, err := cmd.Flags().GetInt("patternID")
		if err != nil {
			log.Fatal(err)
		}
		if patternID == 0 {
			log.Fatal("patternID must be specified.")
		}

		ph := usecase.NewPatternHandler(mysql_client.NewMySQLClient(dbName, user))
		defer ph.Close()
		ph.DeletePattern(patternID)
	},
}

func init() {
	patternCmd.AddCommand(patternDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	patternDeleteCmd.Flags().Int("patternID", 0, "Pattern ID")

	patternDeleteCmd.MarkFlagRequired("patternID")
}
