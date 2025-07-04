package cmd

import (
	"kakeibodb/internal/repository/mysql"
	"kakeibodb/internal/usecase"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// removeTagCmd represents the removeTag command
var patternRemoveTagCmd = &cobra.Command{
	Use:   "removeTag",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		patternID, err := cmd.Flags().GetInt64("patternID")
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
		ptmRepo := mysql.NewPatternTagMapRepository(db)
		ptmUC := usecase.NewPatternTagMapUseCase(ptmRepo)
		err = ptmUC.RemoveTag(patternID, tagName)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	patternCmd.AddCommand(patternRemoveTagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeTagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeTagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	patternRemoveTagCmd.Flags().Int64("patternID", 0, "Pattern ID")
	patternRemoveTagCmd.Flags().StringP("tagName", "t", "", "Tag Name")

	patternRemoveTagCmd.MarkFlagRequired("patternID")
	patternRemoveTagCmd.MarkFlagRequired("tagName")
}
