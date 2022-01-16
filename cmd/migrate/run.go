package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs migrations in the current project",
	Run: func(cmd *cobra.Command, args []string) {
		c := exec.Command("go", "run", "./migrations/.")
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		if err := c.Run(); err != nil {
			fmt.Fprintln(c.Stderr, err)
		}
	},
}
