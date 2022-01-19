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
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var dir string
		if len(args) == 0 {
			dir = "up"
		} else {
			switch args[0] {
			case "up":
				dir = "up"
			case "down":
				dir = "down"
			default:
				dir = "up"
			}
		}

		c := exec.Command("go", "run", "./migrations/.", dir)
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		if err := c.Run(); err != nil {
			fmt.Fprintln(c.Stderr, err)
		}
	},
}
