package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/dotenv"
	"vaultpull/internal/vault"
)

var overwrite bool

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync secrets from Vault into a local .env file",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		client, err := vault.NewClient(cfg.VaultAddr, cfg.Token)
		if err != nil {
			return fmt.Errorf("creating vault client: %w", err)
		}

		secrets, err := client.ReadSecrets(cfg.SecretPath)
		if err != nil {
			return fmt.Errorf("reading secrets from vault: %w", err)
		}

		if len(secrets) == 0 {
			log.Printf("no secrets found at path %s, nothing to write", cfg.SecretPath)
			return nil
		}

		w := dotenv.NewWriter(cfg.OutputFile, overwrite)
		if err := w.Write(secrets); err != nil {
			return fmt.Errorf("writing .env file: %w", err)
		}

		log.Printf("synced %d secret(s) to %s", len(secrets), cfg.OutputFile)
		return nil
	},
}

func init() {
	syncCmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys in .env file")
	rootCmd.AddCommand(syncCmd)
}
