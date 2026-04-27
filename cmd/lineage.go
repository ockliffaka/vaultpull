package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/vaultpull/vaultpull/internal/dotenv"
)

var lineageKey string

var lineageCmd = &cobra.Command{
	Use:   "lineage [env-file]",
	Short: "Show mutation history for a .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		envPath := args[0]

		lin, err := dotenv.LoadLineage(envPath)
		if err != nil {
			return fmt.Errorf("failed to load lineage: %w", err)
		}

		var entries []dotenv.LineageEntry
		if lineageKey != "" {
			entries = lin.KeyHistory(lineageKey)
		} else {
			entries = lin.Entries
		}

		if len(entries) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No lineage records found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TIMESTAMP\tKEY\tOP\tSOURCE\tOLD\tNEW")
		for _, e := range entries {
			old := e.OldValue
			if old == "" {
				old = "-"
			}
			newVal := e.NewValue
			if newVal == "" {
				newVal = "-"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				e.Timestamp.Format("2006-01-02T15:04:05Z"),
				e.Key,
				e.Operation,
				e.Source,
				old,
				newVal,
			)
		}
		return w.Flush()
	},
}

func init() {
	lineageCmd.Flags().StringVar(&lineageKey, "key", "", "Filter history to a single secret key")
	rootCmd.AddCommand(lineageCmd)
}
