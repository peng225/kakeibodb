package cmd

import (
	"github.com/spf13/cobra"
)

// patternCmd represents the pattern command
var patternCmd = &cobra.Command{
	Use:   "pattern",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: nil,
}

func init() {
	rootCmd.AddCommand(patternCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// patternCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// patternCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
