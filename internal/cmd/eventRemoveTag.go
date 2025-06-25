package cmd

import (
	"log"

	"kakeibodb/internal/mysql_client"
	"kakeibodb/internal/usecase"

	"github.com/spf13/cobra"
)

// eventRemoveTagCmd represents the removeTag command
var eventRemoveTagCmd = &cobra.Command{
	Use:   "removeTag",
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
		tagName, err := cmd.Flags().GetString("tagName")
		if err != nil {
			log.Fatal(err)
		}
		if eventID == 0 && tagName == "" {
			log.Fatal("both eventID and tagName must be specified.")
		}

		eh := usecase.NewEventHandler(mysql_client.NewMySQLClient(dbName, dbPort, user))
		defer eh.Close()
		eh.RemoveTag(eventID, tagName)
	},
}

func init() {
	eventCmd.AddCommand(eventRemoveTagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeTagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeTagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	eventRemoveTagCmd.Flags().Int("eventID", 0, "Event ID")
	eventRemoveTagCmd.Flags().StringP("tagName", "t", "", "Tag Name")
}
