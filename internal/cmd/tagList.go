package cmd

import (
	"kakeibodb/internal/mysql_client"
	"kakeibodb/internal/usecase"

	"github.com/spf13/cobra"
)

// tagListCmd represents the tagList command
var tagListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		lh := usecase.NewListHandler(mysql_client.NewMySQLClient(dbName, dbPort, user))
		defer lh.Close()
		lh.ListTag()
	},
}

func init() {
	tagCmd.AddCommand(tagListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
