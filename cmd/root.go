package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "vaultpull",
	Short: "Sync secrets from HashiCorp Vault into local .env files",
	Long: `vaultpull fetches secrets from a HashiCorp Vault KV store
and writes them as key=value pairs into a local .env file,
keeping your development environment in sync safely.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"config file (default: .vaultpull.yaml in current dir or $HOME)",
	)
}
