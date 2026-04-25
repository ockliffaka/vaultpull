package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpull/internal/dotenv"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environment contexts for .env files",
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available environment contexts in the working directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("dir")
		if dir == "" {
			var err error
			dir, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve working dir: %w", err)
			}
		}

		contexts, err := dotenv.ListEnvContexts(dir)
		if err != nil {
			return err
		}

		if len(contexts) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No environment contexts found.")
			return nil
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Available environment contexts:")
		for _, ctx := range contexts {
			path := dotenv.EnvContextPath(dir, ctx)
			fmt.Fprintf(cmd.OutOrStdout(), "  %-16s → %s\n", ctx, path)
		}
		return nil
	},
}

var envShowCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show the resolved .env path for a given context",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("dir")
		if dir == "" {
			var err error
			dir, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve working dir: %w", err)
			}
		}
		name := ""
		if len(args) > 0 {
			name = strings.TrimSpace(args[0])
		}
		ctx := dotenv.ResolveEnvContext(dir, name)
		fmt.Fprintf(cmd.OutOrStdout(), "Context: %s\nPath:    %s\n", ctx.Name, ctx.Path())
		return nil
	},
}

func init() {
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envShowCmd)
	envCmd.PersistentFlags().String("dir", "", "Base directory to scan (defaults to current directory)")
	rootCmd.AddCommand(envCmd)
}
