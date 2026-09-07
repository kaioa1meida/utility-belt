package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd creates the base utility-belt command.
func NewRootCmd(stdout, stderr io.Writer) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "utility-belt",
		Short: "A collection of handy developer utilities for everyday tasks",
		Long: `utility-belt is a fast, offline, and stateless CLI tool providing
handy developer utilities for everyday workflows.`,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	if stdout != nil {
		rootCmd.SetOut(stdout)
	}
	if stderr != nil {
		rootCmd.SetErr(stderr)
	}

	rootCmd.AddCommand(NewGenerateCmd())

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// It is called by main.main().
func Execute() {
	rootCmd := NewRootCmd(os.Stdout, os.Stderr)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

