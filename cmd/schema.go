package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/dotenv"
)

var schemaFile string

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Validate a .env file against a JSON schema of expected keys",
	Long: `Reads a JSON schema file defining required and optional keys (with optional
regex patterns) and validates the current .env file against it.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		envPath, _ := cmd.Flags().GetString("env-file")

		raw, err := os.ReadFile(schemaFile)
		if err != nil {
			return fmt.Errorf("reading schema file: %w", err)
		}

		var rules []dotenv.SchemaRule
		if err := json.Unmarshal(raw, &rules); err != nil {
			return fmt.Errorf("parsing schema file: %w", err)
		}

		secrets, err := dotenv.ParseFile(envPath)
		if err != nil {
			return fmt.Errorf("reading env file: %w", err)
		}

		result := dotenv.ValidateSchema(secrets, rules)
		fmt.Println(result.Summary())

		if !result.OK() {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	schemaCmd.Flags().StringVar(&schemaFile, "schema", ".vault-schema.json", "path to JSON schema file")
	schemaCmd.Flags().String("env-file", ".env", "path to .env file to validate")
	rootCmd.AddCommand(schemaCmd)
}
