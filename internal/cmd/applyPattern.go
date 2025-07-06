package cmd

import (
	"log/slog"
	"os"

	"kakeibodb/internal/model"
	"kakeibodb/internal/repository/mysql"
	"kakeibodb/internal/usecase"

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
		strFrom, err := cmd.Flags().GetString("from")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		strTo, err := cmd.Flags().GetString("to")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		from, err := model.ParseDate(strFrom)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		to, err := model.ParseDate(strTo)
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
		eventRepo := mysql.NewEventRepository(db)
		etmRepo := mysql.NewEventTagMapRepository(db)
		patternRepo := mysql.NewPatternRepository(db)
		apUC := usecase.NewApplyPatternUseCase(eventRepo, etmRepo, patternRepo)
		err = apUC.ApplyPattern(*from, *to)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
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
	applyPatternCmd.Flags().String("from", "2018-01-01", "the beginning of time range")
	applyPatternCmd.Flags().String("to", "2100-12-31", "the end of time range")
}
