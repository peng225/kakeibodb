package cmd

import (
	"context"
	"log/slog"
	"os"

	"kakeibodb/internal/repository/mysql"
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
			slog.Error(err.Error())
			os.Exit(1)
		}
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		credit, err := cmd.Flags().GetBool("credit")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		parentEventID, err := cmd.Flags().GetInt64("parentEventID")
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
		eventUC := usecase.NewEventUseCase(eventRepo, tx)

		ctx := context.Background()
		if credit {
			if file == "" {
				slog.Error("File path must be specified for credit mode.")
				os.Exit(1)
			}
			if parentEventID < 0 {
				slog.Error("Invalid argument.", "parentEventID", parentEventID)
				os.Exit(1)
			}
			err = eventUC.LoadCreditFromFile(ctx, file, parentEventID)
		} else {
			if file != "" {
				err = eventUC.LoadFromFile(ctx, file)
			} else {
				err = eventUC.LoadFromDir(ctx, dir)
			}
		}
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
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
	loadCmd.Flags().Bool("credit", false, "Load credit card event data")
	loadCmd.Flags().Int64("parentEventID", -1, "The parent event ID related to the credit events to be loaded")

	loadCmd.MarkFlagsMutuallyExclusive("file", "dir")
	loadCmd.MarkFlagsOneRequired("file", "dir")
}
