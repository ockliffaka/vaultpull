package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/dotenv"
)

var ttlCmd = &cobra.Command{
	Use:   "ttl",
	Short: "Show the TTL status of a local .env file",
	RunE: func(cmd *cobra.Command, args []string) error {
		envFile, _ := cmd.Flags().GetString("env-file")
		maxAge, _ := cmd.Flags().GetDuration("max-age")
		warnAge, _ := cmd.Flags().GetDuration("warn-age")

		policy := dotenv.TTLPolicy{
			MaxAge:  maxAge,
			WarnAge: warnAge,
		}

		summary, err := dotenv.TTLSummary(envFile, policy)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return err
		}

		fmt.Println(summary)
		return nil
	},
}

func init() {
	ttlCmd.Flags().String("env-file", ".env", "path to the .env file")
	ttlCmd.Flags().Duration("max-age", 24*60*60*1000000000, "maximum secret age before expiry")
	ttlCmd.Flags().Duration("warn-age", 20*60*60*1000000000, "age threshold for warning")
	rootCmd.AddCommand(ttlCmd)
}
