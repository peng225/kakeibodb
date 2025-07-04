package cmd

import (
	"kakeibodb/internal/repository/mysql"
	"kakeibodb/internal/usecase"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// eventRemoveTagCmd represents the removeTag command
var eventRemoveTagCmd = &cobra.Command{
	Use:   "removeTag",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		eventID, err := cmd.Flags().GetInt64("eventID")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		tagName, err := cmd.Flags().GetString("tagName")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		db, err := OpenDB(dbName, dbPort, user)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		defer db.Close()
		etmRepo := mysql.NewEventTagMapRepository(db)
		etmUC := usecase.NewEventTagMapUseCase(etmRepo)
		err = etmUC.RemoveTag(eventID, tagName)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	eventCmd.AddCommand(eventRemoveTagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeTagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeTagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	eventRemoveTagCmd.Flags().Int64("eventID", 0, "Event ID")
	eventRemoveTagCmd.Flags().StringP("tagName", "t", "", "Tag Name")

	eventRemoveTagCmd.MarkFlagRequired("eventID")
	eventRemoveTagCmd.MarkFlagRequired("tagName")
}
