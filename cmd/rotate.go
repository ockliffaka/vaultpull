package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/dotenv"
	"vaultpull/internal/vault"
)

var rotateTTL time.Duration

var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate secrets: backup, overwrite, and refresh expiry stamp",
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

		result, err := dotenv.Rotate(cfg.EnvFile, secrets, rotateTTL)
		if err != nil {
			return fmt.Errorf("rotate: %w", err)
		}

		fmt.Printf("✔ Rotated %d keys into %s\n", result.KeysUpdated, result.Path)
		fmt.Printf("  Backup: %s\n", result.BackupPath)
		fmt.Printf("  Expires in: %s\n", rotateTTL)
		return nil
	},
}

func init() {
	rotateCmd.Flags().DurationVar(&rotateTTL, "ttl", 24*time.Hour, "TTL before secrets are considered expired")
	rootCmd.AddCommand(rotateCmd)
}
