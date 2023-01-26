package cmd

import (
	"kakeibodb/mysql_client"
	"kakeibodb/usecase"
	"log"

	"github.com/spf13/cobra"
)

// applyPatternCmd represents the applyPattern command
var applyPatternCmd = &cobra.Command{
	Use:   "applyPattern",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		from, err := cmd.Flags().GetString("from")
		if err != nil {
			log.Fatal(err)
		}
		to, err := cmd.Flags().GetString("to")
		if err != nil {
			log.Fatal(err)
		}

		eh := usecase.NewEventHandler(mysql_client.NewMySQLClient(dbName, user))
		defer eh.Close()
		eh.ApplyPattern(from, to)
	},
}

func init() {
	eventCmd.AddCommand(applyPatternCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyPatternCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyPatternCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	applyPatternCmd.Flags().StringP("from", "", "2018-01-01", "the beginning of time range")
	applyPatternCmd.Flags().StringP("to", "", "2100-12-31", "the end of time range")
}
