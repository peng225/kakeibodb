package cmd

import (
	"context"
	"fmt"
	"kakeibodb/internal/model"
	"kakeibodb/internal/repository/mysql"
	"kakeibodb/internal/usecase"
	"log/slog"
	"os"

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
			slog.Info("Env var detected.", envSplitBaseTagName, splitBaseTagName)
		}
		// The env can be empty string. That is also OK.
	},
	Run: func(cmd *cobra.Command, args []string) {
		eventIDs, err := cmd.Flags().GetInt64Slice("eventIDs")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		if len(eventIDs) == 0 && splitBaseTagName == "" {
			slog.Error(
				fmt.Sprintf("Either --eventIDs flag or %s env should be set.",
					envSplitBaseTagName),
			)
			os.Exit(1)
		}
		strDate, err := cmd.Flags().GetString("date")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		money, err := cmd.Flags().GetInt32("money")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		desc, err := cmd.Flags().GetString("desc")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		date, err := model.ParseDate(strDate)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		db, err := OpenDB(dbName, dbPort, user)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		defer db.Close()
		eventRepo := mysql.NewEventRepository(db)
		tx := mysql.NewTransaction(db)
		eventUC := usecase.NewEventUseCase(eventRepo, nil, tx)
		ctx := context.Background()
		err = eventUC.Split(ctx, eventIDs, splitBaseTagName, *date, money, desc)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
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
	splitCmd.Flags().Int64Slice("eventIDs", []int64{}, "The list of event ID(s) to be split")
	splitCmd.Flags().String("date", "", "Date of the new event (YYYY-MM-DD or YYYY/MM/DD)")
	splitCmd.Flags().Int32("money", -1, "Money of the new event")
	splitCmd.Flags().String("desc", "", "Description of the new event")

	splitCmd.MarkFlagRequired("date")
	splitCmd.MarkFlagRequired("money")
	splitCmd.MarkFlagRequired("desc")
}
