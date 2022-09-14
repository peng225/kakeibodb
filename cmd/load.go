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
			log.Fatal(err)
		}
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			log.Fatal(err)
		}
		leh := usecase.NewLoadEventHandler(mysql_client.NewMySQLClient())
		if file == "" && dir == "" {
			log.Fatal("either file or dir must be specified.")
		} else if file != "" && dir != "" {
			log.Fatal("both file and dir cannot be specified.")
		} else if file != "" {
			leh.LoadEventFromFile(file)
		} else {
			leh.LoadEventFromDir(dir)
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
}
