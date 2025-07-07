/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log/slog"
	"os"

	"kakeibodb/internal/model"
	"kakeibodb/internal/presenter/console"
	"kakeibodb/internal/repository/mysql"
	"kakeibodb/internal/usecase"

	"github.com/spf13/cobra"
)

// tsCmd represents the ts command
var tsCmd = &cobra.Command{
	Use:   "ts",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cf, err := parseCommonFlags(cmd)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		interval, err := cmd.Flags().GetInt("interval")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		window, err := cmd.Flags().GetInt("window")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		top, err := cmd.Flags().GetInt("top")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		from, err := model.ParseDate(cf.from)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		to, err := model.ParseDate(cf.to)
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
		analysysPresenter := console.NewAnalysisPresenter()
		analysisUC := usecase.NewAnalysisUseCase(eventRepo, analysysPresenter)
		ctx := context.Background()
		err = analysisUC.TimeSeries(ctx, *from, *to, interval, window, top)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	analysisCmd.AddCommand(tsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	tsCmd.Flags().Int("window", 3, "time window (month) for time series calculation")
	tsCmd.Flags().Int("interval", 1, "interval (month) for time series calculation")
	tsCmd.Flags().Int("top", 10, "the number of results for time series calculation")
}
