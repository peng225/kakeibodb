/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"kakeibodb/internal/mysql_client"
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
			log.Fatal(err)
		}

		interval, err := cmd.Flags().GetInt("interval")
		if err != nil {
			log.Fatal(err)
		}
		window, err := cmd.Flags().GetInt("window")
		if err != nil {
			log.Fatal(err)
		}
		top, err := cmd.Flags().GetInt("top")
		if err != nil {
			log.Fatal(err)
		}

		ah := usecase.NewAnalysisHandler(mysql_client.NewMySQLClient(dbName, dbPort, user))
		ah.TimeSeries(cf.from, cf.to, interval, window, top)
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
	tsCmd.Flags().IntP("window", "", 3, "time window (month) for time series calculation")
	tsCmd.Flags().IntP("interval", "", 1, "interval (month) for time series calculation")
	tsCmd.Flags().IntP("top", "", 10, "the number of results for time series calculation")
}
