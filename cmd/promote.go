package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/dotenv"
)

var (
	promoteOverwrite bool
	promoteDryRun    bool
	promoteKeys      []string
)

var promoteCmd = &cobra.Command{
	Use:   "promote <source-env> <target-env>",
	Short: "Promote secrets from one environment to another",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcLabel := args[0]
		dstLabel := args[1]

		srcPath, err := dotenv.ResolveEnvContext(srcLabel)
		if err != nil {
			return fmt.Errorf("resolve source env %q: %w", srcLabel, err)
		}
		dstPath, err := dotenv.ResolveEnvContext(dstLabel)
		if err != nil {
			return fmt.Errorf("resolve target env %q: %w", dstLabel, err)
		}

		src, err := dotenv.LoadSnapshot(srcPath)
		if err != nil || src == nil {
			return fmt.Errorf("load source env %q: %w", srcPath, err)
		}
		dst, err := dotenv.LoadSnapshot(dstPath)
		if err != nil {
			return fmt.Errorf("load target env %q: %w", dstPath, err)
		}
		if dst == nil {
			dst = map[string]string{}
		}

		opts := dotenv.PromoteOptions{
			Overwrite: promoteOverwrite,
			DryRun:    promoteDryRun,
			Keys:      promoteKeys,
		}

		out, result := dotenv.Promote(src, dst, srcLabel, dstLabel, opts)
		fmt.Fprintln(os.Stdout, result.Summary())

		if promoteDryRun {
			fmt.Fprintln(os.Stdout, "[dry-run] no changes written")
			return nil
		}

		return dotenv.SaveSnapshot(dstPath, out)
	},
}

func init() {
	promoteCmd.Flags().BoolVar(&promoteOverwrite, "overwrite", false, "Overwrite existing keys in target env")
	promoteCmd.Flags().BoolVar(&promoteDryRun, "dry-run", false, "Preview changes without writing")
	promoteCmd.Flags().StringSliceVar(&promoteKeys, "keys", nil, "Comma-separated list of keys to promote")
	rootCmd.AddCommand(promoteCmd)
}
