package cmd

import (
	"log"

	"kakeibodb/internal/repository/mysql"
	"kakeibodb/internal/usecase"

	"github.com/spf13/cobra"
)

// patternDeleteCmd represents the delete command
var patternDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		patternID, err := cmd.Flags().GetInt64("patternID")
		if err != nil {
			log.Fatal(err)
		}
		if patternID == 0 {
			log.Fatal("patternID must be specified.")
		}

		db, err := OpenDB(dbName, dbPort, user)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		patternRepo := mysql.NewPatternRepository(db)
		patternUC := usecase.NewPatternUseCase(patternRepo)
		patternUC.Delete(patternID)
	},
}

func init() {
	patternCmd.AddCommand(patternDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	patternDeleteCmd.Flags().Int64("patternID", 0, "Pattern ID")

	patternDeleteCmd.MarkFlagRequired("patternID")
}
