package cmd

import (
	"log"

	"kakeibodb/internal/repository/mysql"
	"kakeibodb/internal/usecase"

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
			log.Fatal(err)
		}
		tagNames, err := cmd.Flags().GetStringSlice("tagNames")
		if err != nil {
			log.Fatal(err)
		}
		if patternID == 0 && len(tagNames) == 0 {
			log.Fatal("both patternID and tagNames must be specified.")
		}

		db, err := OpenDB(dbName, dbPort, user)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		ptmRepo := mysql.NewPatternTagMapRepository(db)
		ptmUC := usecase.NewPatternTagMapUseCase(ptmRepo)
		ptmUC.AddTag(patternID, tagNames)
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
