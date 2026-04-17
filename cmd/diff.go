package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/dotenv"
	"vaultpull/internal/vault"
)

var diffShowValues bool

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Preview changes between Vault secrets and local .env file",
	Long:  `Fetches secrets from Vault and shows a diff against the local .env file without writing any changes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}

		client, err := vault.NewClient(cfg)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets, err := client.ReadSecrets(cfg.SecretPath)
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		existing := map[string]string{}
		if _, statErr := os.Stat(cfg.OutputFile); statErr == nil {
			existing, err = dotenv.Parse(cfg.OutputFile)
			if err != nil {
				return fmt.Errorf("parse existing env: %w", err)
			}
		}

		diff := dotenv.Diff(existing, secrets)
		dotenv.PrintPreview(os.Stdout, diff, diffShowValues)
		return nil
	},
}

func init() {
	diffCmd.Flags().BoolVar(&diffShowValues, "show-values", false, "Show actual secret values in diff output")
	rootCmd.AddCommand(diffCmd)
}
