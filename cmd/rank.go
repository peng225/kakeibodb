/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"kakeibodb/mysql_client"
	"kakeibodb/usecase"
	"log"

	"github.com/spf13/cobra"
)

// rankCmd represents the rank command
var rankCmd = &cobra.Command{
	Use:   "rank",
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

		ah := usecase.NewAnalysisHandler(mysql_client.NewMySQLClient(dbName, user))
		ah.Rank(cf.from, cf.to)
	},
}

func init() {
	analysisCmd.AddCommand(rankCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rankCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rankCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
