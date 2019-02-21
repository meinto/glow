package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version of glow",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version: 0.7.0")
	},
}
