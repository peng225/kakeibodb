/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// analysisCmd represents the analysis command
var analysisCmd = &cobra.Command{
	Use:   "analysis",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(analysisCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// analysisCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// analysisCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	analysisCmd.PersistentFlags().String("from", "2018-01-01", "the beginning of time range")
	analysisCmd.PersistentFlags().String("to", "2100-12-31", "the end of time range")
}

type commonFlags struct {
	from string
	to   string
}

func parseCommonFlags(cmd *cobra.Command) (*commonFlags, error) {
	from, err := cmd.Flags().GetString("from")
	if err != nil {
		return nil, err
	}
	to, err := cmd.Flags().GetString("to")
	if err != nil {
		return nil, err
	}
	return &commonFlags{
		from: from,
		to:   to,
	}, nil
}
