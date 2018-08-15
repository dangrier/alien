package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "alien",
	Short: "Alien probes your endpoints",
	Long: `Alien probes your endpoints and filters the result
to determine whether the probe was successful.
Results are exposed as Prometheus metrics.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
