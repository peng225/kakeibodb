package cmd

import (
	"kakeibodb/internal/presenter/console"
	"kakeibodb/internal/repository/mysql"
	"kakeibodb/internal/usecase"
	"log"

	"github.com/spf13/cobra"
)

// patternListCmd represents the list command
var patternListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := OpenDB(dbName, dbPort, user)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		patternRepo := mysql.NewPatternRepository(db)
		patternPresenter := console.NewPatternPresenter()
		patternUC := usecase.NewPatternPresentUseCase(patternRepo, patternPresenter)
		patternUC.List()
	},
}

func init() {
	patternCmd.AddCommand(patternListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
