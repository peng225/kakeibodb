package cmd

import (
	"log"

	"kakeibodb/internal/mysql_client"
	"kakeibodb/internal/usecase"

	"github.com/spf13/cobra"
)

// tagCreateCmd represents the tagCreate command
var tagCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		tagName, err := cmd.Flags().GetString("tagName")
		if err != nil {
			log.Fatal(err)
		}
		if tagName == "" {
			log.Fatal("tag name must be specified.")
		}

		th := usecase.NewTagHandler(mysql_client.NewMySQLClient(dbName, dbPort, user))
		defer th.Close()
		th.CreateTag(tagName)
	},
}

func init() {
	tagCmd.AddCommand(tagCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagCreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	tagCreateCmd.Flags().StringP("tagName", "t", "", "Tag name")
}
