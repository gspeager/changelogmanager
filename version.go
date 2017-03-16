package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the current version of Changelog Manager",
	Long:  "Displays the current version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Changelog Manager v" + version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
