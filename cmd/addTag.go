/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"kakeibodb/mysql_client"
	"kakeibodb/usecase"
	"log"

	"github.com/spf13/cobra"
)

// addTagCmd represents the addTag command
var addTagCmd = &cobra.Command{
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
		tagName, err := cmd.Flags().GetString("tagName")
		if err != nil {
			log.Fatal(err)
		}
		if eventID == 0 && tagName == "" {
			log.Fatal("both eventID and tagName must be specified.")
		}

		eh := usecase.NewEventHandler(mysql_client.NewMySQLClient())
		eh.AddTag(eventID, tagName)
	},
}

func init() {
	eventCmd.AddCommand(addTagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addTagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addTagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addTagCmd.Flags().IntP("eventID", "e", 0, "Event ID")
	addTagCmd.Flags().StringP("tagName", "t", "", "Tag Name")
}
