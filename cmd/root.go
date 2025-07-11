package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goat",
	Short: "A CLI tool to generate Go Application Tmeplate.",
	Long: `goat is a simple command-line interface
to help you quickly bootstrap your Go Applications Template.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func init() {
}
