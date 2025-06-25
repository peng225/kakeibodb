package cmd

import (
	"log"

	"kakeibodb/internal/mysql_client"
	"kakeibodb/internal/usecase"

	"github.com/spf13/cobra"
)

// tagDeleteCmd represents the tagDelete command
var tagDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		tagID, err := cmd.Flags().GetInt("tagID")
		if err != nil {
			log.Fatal(err)
		}
		if tagID == 0 {
			log.Fatal("tagID must be specified.")
		}

		th := usecase.NewTagHandler(mysql_client.NewMySQLClient(dbName, user))
		defer th.Close()
		th.DeleteTag(tagID)
	},
}

func init() {
	tagCmd.AddCommand(tagDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagDeleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagDeleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	tagDeleteCmd.Flags().Int("tagID", 0, "Tag ID")
}
