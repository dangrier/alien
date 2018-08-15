package commands

import (
	"github.com/dangrier/alien/pkg/alien"
	"github.com/dangrier/alien/pkg/probe"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cmdTest = &cobra.Command{
	Use:   "test",
	Short: "test the current build",
	Run: func(cmd *cobra.Command, args []string) {
		doTest()
	},
}

func init() {
	rootCmd.AddCommand(cmdTest)
}

func doTest() {
	a := alien.New()

	p, err := probe.New("https://google.com.au", probe.WithSuccessFilter(probe.FilterResponseCode(200)))
	if err != nil {
		logrus.Fatalf("New probe: %v", err)
	}

	err = a.AddProbe(p)
	if err != nil {
		logrus.Fatalf("Add probe: %v", err)
	}

	a.ListenForTermination()
	a.Run()
}
