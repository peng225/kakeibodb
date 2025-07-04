package cmd

import (
	"kakeibodb/internal/model"
	"kakeibodb/internal/presenter/console"
	"kakeibodb/internal/repository/mysql"
	"kakeibodb/internal/usecase"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// eventListCmd represents the eventList command
var eventListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		tagNames, err := cmd.Flags().GetStringSlice("tags")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		fromStr, err := cmd.Flags().GetString("from")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		from, err := model.ParseDate(fromStr)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		toStr, err := cmd.Flags().GetString("to")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		to, err := model.ParseDate(toStr)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		lastDays, err := cmd.Flags().GetInt("last")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		if lastDays >= 0 {
			tmpFrom := time.Now().AddDate(0, 0, -lastDays)
			from = &tmpFrom
		}

		db, err := OpenDB(dbName, dbPort, user)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		defer db.Close()
		eventRepo := mysql.NewEventRepository(db)
		eventPresenter := console.NewEventPresenter()
		eventPresentUC := usecase.NewEventPresentUseCase(eventRepo, eventPresenter)
		if all {
			err = eventPresentUC.PresentAll(tagNames, from, to)
		} else {
			err = eventPresentUC.PresentOutcomes(tagNames, from, to)
		}
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	eventCmd.AddCommand(eventListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// eventListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// eventListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	eventListCmd.Flags().StringSlice("tags", nil, `tag list (eg. "foo", "foo,var" etc.)`)
	eventListCmd.Flags().String("from", "2018-01-01", "the beginning of time range")
	eventListCmd.Flags().String("to", "2100-12-31", "the end of time range")
	eventListCmd.Flags().Int("last", -1, "show the events of last X days")
	eventListCmd.Flags().BoolP("all", "a", false, "show all events")

	eventListCmd.MarkFlagsMutuallyExclusive("from", "last")
	eventListCmd.MarkFlagsMutuallyExclusive("to", "last")
}
