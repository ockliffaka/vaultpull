package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/dotenv"
)

var pinCmd = &cobra.Command{
	Use:   "pin",
	Short: "Pin the current .env secrets under a named label",
	Long: `Save a named snapshot (pin) of the current .env file's secrets.
Pins can be loaded later for diffing or rollback purposes.

Example:
  vaultpull pin --file .env --label release-1.2.0`,
	RunE: func(cmd *cobra.Command, args []string) error {
		envFile, _ := cmd.Flags().GetString("file")
		label, _ := cmd.Flags().GetString("label")
		delete, _ := cmd.Flags().GetBool("delete")

		if label == "" {
			return fmt.Errorf("--label is required")
		}

		if delete {
			if err := dotenv.DeletePin(envFile, label); err != nil {
				return fmt.Errorf("failed to delete pin: %w", err)
			}
			fmt.Fprintf(os.Stdout, "Pin %q deleted.\n", label)
			return nil
		}

		writer, err := dotenv.NewWriter(envFile)
		if err != nil {
			return fmt.Errorf("failed to read env file: %w", err)
		}
		secrets := writer.Current()

		if err := dotenv.Pin(envFile, label, secrets); err != nil {
			return fmt.Errorf("failed to create pin: %w", err)
		}
		fmt.Fprintf(os.Stdout, "Pinned %d secrets under label %q.\n", len(secrets), label)
		return nil
	},
}

func init() {
	pinCmd.Flags().String("file", ".env", "Path to the .env file to pin")
	pinCmd.Flags().String("label", "", "Label name for the pin (required)")
	pinCmd.Flags().Bool("delete", false, "Delete the named pin instead of creating one")
	rootCmd.AddCommand(pinCmd)
}
