package cmd

import (
	"log"

	"kakeibodb/internal/mysql_client"
	"kakeibodb/internal/usecase"

	"github.com/spf13/cobra"
)

// eventAddTagCmd represents the addTag command
var eventAddTagCmd = &cobra.Command{
	Use:   "addTag",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		eventID, err := cmd.Flags().GetInt("eventID")
		if err != nil {
			log.Fatal(err)
		}
		tagNames, err := cmd.Flags().GetStringSlice("tagNames")
		if err != nil {
			log.Fatal(err)
		}
		if eventID == 0 && len(tagNames) == 0 {
			log.Fatal("both eventID and tagNames must be specified.")
		}

		eh := usecase.NewEventHandler(mysql_client.NewMySQLClient(dbName, dbPort, user))
		defer eh.Close()
		eh.AddTag(eventID, tagNames)
	},
}

func init() {
	eventCmd.AddCommand(eventAddTagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addTagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addTagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	eventAddTagCmd.Flags().Int("eventID", 0, "credit == false: Event ID, credit == true: Credit card event ID")
	eventAddTagCmd.Flags().StringSlice("tagNames", nil, "Tag Names")
}
