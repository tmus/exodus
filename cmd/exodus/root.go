package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "exodus",
	Short: "Exodus is a database migration runner and DSL for Go",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(initCmd, createCmd, runCmd)
}
