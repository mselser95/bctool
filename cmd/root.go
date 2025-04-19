/*
Copyright © 2025
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bctool",
	Short: "Blockchain information tool",
	Long: `bctool is a CLI application that provides access to various blockchain-related utilities.
You can fetch prices, inspect addresses, or analyze data across different chains.

Examples:
  bctool prices BTC ETH TRC
`,
}

// Execute runs the root command and its subcommands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add global flags here if needed in the future
}
