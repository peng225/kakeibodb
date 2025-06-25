package cmd

import (
	"fmt"
	"kakeibodb/internal/db_client"
	"log"
	"os"
	"time"

	"kakeibodb/internal/mysql_client"
	"kakeibodb/internal/usecase"

	"github.com/spf13/cobra"
)

const envSplitBaseTagName = "KAKEIBODB_SPLIT_BASE_TAG_NAME"

var splitBaseTagName string

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		splitBaseTagName = os.Getenv(envSplitBaseTagName)
		if splitBaseTagName != "" {
			log.Printf("%s: %s\n", envSplitBaseTagName, splitBaseTagName)
		}
		// The env can be empty string. That is also OK.
	},
	Run: func(cmd *cobra.Command, args []string) {
		eventID, err := cmd.Flags().GetInt("eventID")
		if err != nil {
			log.Fatal(err)
		}
		if eventID == -1 && splitBaseTagName == "" {
			log.Fatalf("Either --eventID flag or %s env should be set.",
				envSplitBaseTagName)
		}
		date, err := cmd.Flags().GetString("date")
		if err != nil {
			log.Fatal(err)
		}
		money, err := cmd.Flags().GetInt("money")
		if err != nil {
			log.Fatal(err)
		}
		desc, err := cmd.Flags().GetString("desc")
		if err != nil {
			log.Fatal(err)
		}

		layouts := []string{"2006/01/02", "2006-01-02"}
		for _, layout := range layouts {
			_, err = time.Parse(layout, date)
			if err == nil {
				break
			}
		}
		if err != nil {
			log.Fatal(err)
		}
		if len([]rune(desc)) >= db_client.EventDescLength {
			desc = string([]rune(desc)[0:db_client.EventDescLength])
		}

		eh := usecase.NewEventHandler(mysql_client.NewMySQLClient(dbName, user))
		defer eh.Close()

		if eventID == -1 {
			eventID, err = eh.GetEventIDFromSplitBaseTag(splitBaseTagName, date)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Auto detected eventID: %d\n", eventID)
		}
		eh.Split(eventID, date, money, desc)
	},
}

func init() {
	eventCmd.AddCommand(splitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// splitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// splitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	splitCmd.Flags().Int("eventID", -1, "The event ID to be split")
	splitCmd.Flags().String("date", "", "Date of the new event (YYYY-MM-DD or YYYY/MM/DD)")
	splitCmd.Flags().Int("money", -1, "Money of the new event")
	splitCmd.Flags().String("desc", "", "Description of the new event")

	splitCmd.MarkFlagRequired("date")
	splitCmd.MarkFlagRequired("money")
	splitCmd.MarkFlagRequired("desc")
}
