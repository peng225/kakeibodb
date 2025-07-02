package cmd

import (
	"log"

	"kakeibodb/internal/repository/mysql"
	"kakeibodb/internal/usecase"

	"github.com/spf13/cobra"
)

// patternCreateCmd represents the create command
var patternCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		key, err := cmd.Flags().GetString("key")
		if err != nil {
			log.Fatal(err)
		}
		if key == "" {
			log.Fatal("Key string must be specified.")
		}

		db, err := OpenDB(dbName, dbPort, user)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		patternRepo := mysql.NewPatternRepository(db)
		patternUC := usecase.NewPatternUseCase(patternRepo)
		patternUC.Create(key)
	},
}

func init() {
	patternCmd.AddCommand(patternCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	patternCreateCmd.Flags().StringP("key", "k", "", "Key string")

	patternCreateCmd.MarkFlagRequired("key")
}
