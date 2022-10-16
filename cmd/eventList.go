/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"kakeibodb/mysql_client"
	"kakeibodb/usecase"
	"log"
	"strings"

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
		tags, err := cmd.Flags().GetString("tags")
		if err != nil {
			log.Fatal(err)
		}
		if strings.Contains(tags, "&") && strings.Contains(tags, "|") {
			log.Fatal(`tags cannot contain both "&" and "|".`)
		}
		from, err := cmd.Flags().GetString("from")
		if err != nil {
			log.Fatal(err)
		}
		to, err := cmd.Flags().GetString("to")
		if err != nil {
			log.Fatal(err)
		}
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			log.Fatal(err)
		}

		lh := usecase.NewListHandler(mysql_client.NewMySQLClient())
		if all {
			lh.ListAllEvent(tags, from, to)
		} else {
			lh.ListPaymentEvent(tags, from, to)
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
	eventListCmd.Flags().StringP("tags", "", "", `tag list (eg. "foo", "foo&var", "foo|var" etc.)`)
	eventListCmd.Flags().StringP("from", "", "2018-01-01", "the beginning of time range")
	eventListCmd.Flags().StringP("to", "", "2100-12-31", "the end of time range")
	eventListCmd.Flags().BoolP("all", "a", false, "show all events")
}
