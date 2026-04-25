package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/dotenv"
	"github.com/your-org/vaultpull/internal/vault"
)

var templateCmd = &cobra.Command{
	Use:   "template <template-file>",
	Short: "Render a .env template file using secrets from Vault",
	Long: `Reads a template file containing \${KEY} placeholders and replaces them
with secrets fetched from HashiCorp Vault. The result is written to --out or
printed to stdout.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tmplPath := args[0]
		outPath, _ := cmd.Flags().GetString("out")
		strict, _ := cmd.Flags().GetBool("strict")

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
			return fmt.Errorf("vault read: %w", err)
		}

		opts := dotenv.DefaultTemplateOptions()
		opts.Strict = strict

		if err := dotenv.RenderFile(tmplPath, secrets, opts, outPath); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return err
		}
		if outPath != "" {
			fmt.Fprintf(os.Stderr, "rendered template → %s\n", outPath)
		}
		return nil
	},
}

func init() {
	templateCmd.Flags().StringP("out", "o", "", "output file path (default: stdout)")
	templateCmd.Flags().Bool("strict", false, "fail if any placeholder cannot be resolved")
	rootCmd.AddCommand(templateCmd)
}
