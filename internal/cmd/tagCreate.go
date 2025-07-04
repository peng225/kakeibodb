package cmd

import (
	"kakeibodb/internal/repository/mysql"
	"kakeibodb/internal/usecase"
	"log/slog"
	"os"

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
			slog.Error(err.Error())
			os.Exit(1)
		}

		db, err := OpenDB(dbName, dbPort, user)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		defer db.Close()
		tagRepo := mysql.NewTagRepository(db)
		tagUC := usecase.NewTagUseCase(tagRepo)
		err = tagUC.Create(tagName)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
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

	tagCreateCmd.MarkFlagRequired("tagName")
}
