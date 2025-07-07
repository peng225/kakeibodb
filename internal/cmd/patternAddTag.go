package cmd

import (
	"context"
	"kakeibodb/internal/repository/mysql"
	"kakeibodb/internal/usecase"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// addTagCmd represents the addTag command
var patternAddTagCmd = &cobra.Command{
	Use:   "addTag",
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
		tagNames, err := cmd.Flags().GetStringSlice("tagNames")
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
		ctx := context.Background()
		err = ptmUC.AddTag(ctx, patternID, tagNames)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	patternCmd.AddCommand(patternAddTagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addTagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addTagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	patternAddTagCmd.Flags().Int64("patternID", 0, "Pattern ID")
	patternAddTagCmd.Flags().StringSlice("tagNames", nil, "Tag Names")

	patternAddTagCmd.MarkFlagRequired("patternID")
	patternAddTagCmd.MarkFlagRequired("tagNames")
}
