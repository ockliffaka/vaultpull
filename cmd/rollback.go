package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/dotenv"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Restore the most recent backup of the .env file",
	RunE: func(cmd *cobra.Command, args []string) error {
		envFile, _ := cmd.Flags().GetString("env-file")

		list, _ := cmd.Flags().GetBool("list")
		if list {
			backups, err := dotenv.ListBackups(envFile)
			if err != nil {
				return err
			}
			if len(backups) == 0 {
				fmt.Println("No backups found.")
				return nil
			}
			for _, b := range backups {
				fmt.Println(b)
			}
			return nil
		}

		restored, err := dotenv.Rollback(envFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "rollback failed: %v\n", err)
			return err
		}
		fmt.Printf("Rolled back %s from %s\n", envFile, restored)
		return nil
	},
}

func init() {
	rollbackCmd.Flags().String("env-file", ".env", "Path to the .env file")
	rollbackCmd.Flags().Bool("list", false, "List available backups without restoring")
	rootCmd.AddCommand(rollbackCmd)
}
