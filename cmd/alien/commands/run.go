package commands

import (
	"github.com/dangrier/alien/pkg/alien"
	"github.com/dangrier/alien/pkg/probe"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "run a single probe with default settings and looking for a HTTP 200 response",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		run(args)
	},
}

func init() {
	rootCmd.AddCommand(cmdRun)
}

func run(endpoints []string) {
	a := alien.New()

	for _, ep := range endpoints {
		p, err := probe.New(ep, probe.WithSuccessFilter(probe.FilterResponseCode(200)))
		if err != nil {
			logrus.Fatalf("New probe: %v", err)
		}

		err = a.AddProbe(p)
		if err != nil {
			logrus.Fatalf("Add probe: %v", err)
		}
	}

	a.Run()
}
