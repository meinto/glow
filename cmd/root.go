package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var author string

var rootCmd = &cobra.Command{
	Use:   "glow",
	Short: "small tool to adapt git-flow for gitlab",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&author, "author", "a", "test", "name of the author")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
