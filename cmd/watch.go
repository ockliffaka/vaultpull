package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"vaultpull/internal/dotenv"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch .env file and refresh secrets when TTL expires",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("file")
		intervalSec, _ := cmd.Flags().GetInt("interval")

		policy := dotenv.DefaultTTLPolicy()
		opts := dotenv.WatchOptions{
			Interval:  time.Duration(intervalSec) * time.Second,
			MaxCycles: 0,
			OnRefresh: func(p string) error {
				fmt.Fprintf(os.Stdout, "[watch] secrets expired, refreshing %s\n", p)
				// In a full implementation this would re-run the sync pipeline.
				return nil
			},
		}

		stop := make(chan struct{})
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			fmt.Fprintln(os.Stdout, "\n[watch] stopping")
			close(stop)
		}()

		fmt.Fprintf(os.Stdout, "[watch] watching %s every %ds\n", path, intervalSec)
		dotenv.Watch(path, policy, opts, stop)
		return nil
	},
}

func init() {
	watchCmd.Flags().String("file", ".env", "path to .env file")
	watchCmd.Flags().Int("interval", 60, "poll interval in seconds")
	rootCmd.AddCommand(watchCmd)
}
