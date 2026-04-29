package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/dotenv"
)

var hookCmd = &cobra.Command{
	Use:   "hook <event> <command> [args...]",
	Short: "Run lifecycle hooks manually for a given event",
	Long: `Manually trigger all registered hooks for a lifecycle event.

Supported events: pre-sync, post-sync, pre-write, post-write

Example:
  vaultpull hook post-sync
  vaultpull hook post-write -- ./scripts/notify.sh`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		event := dotenv.HookEvent(args[0])
		switch event {
		case dotenv.HookPreSync, dotenv.HookPostSync, dotenv.HookPreWrite, dotenv.HookPostWrite:
			// valid
		default:
			return fmt.Errorf("unknown event %q; valid events: pre-sync, post-sync, pre-write, post-write", event)
		}

		var hooks []dotenv.Hook
		if len(args) >= 2 {
			hooks = append(hooks, dotenv.Hook{
				Event:   event,
				Command: args[1],
				Args:    args[2:],
			})
		}

		stopOnError, _ := cmd.Flags().GetBool("stop-on-error")
		results, err := dotenv.RunHooks(hooks, event, stopOnError)

		for _, r := range results {
			if r.Output != "" {
				fmt.Fprintf(os.Stdout, "[%s] %s\n%s\n", r.Event, r.Command, r.Output)
			}
			if r.Err != nil {
				fmt.Fprintf(os.Stderr, "hook error: %v\n", r.Err)
			}
		}

		fmt.Println(dotenv.HookSummary(results))

		return err
	},
}

func init() {
	hookCmd.Flags().Bool("stop-on-error", false, "stop executing hooks after the first failure")
	rootCmd.AddCommand(hookCmd)
}
