package leszy

import (
	"github.com/SamsungSLAV/slav/logger"
	"github.com/spf13/cobra"
)

type rootCmd struct {
	BaseCmd
}

func NewRootCmd(c *Clients) *rootCmd {
	return &rootCmd{BaseCmd{
		Clients: c,
		Command: &cobra.Command{
			Use:   "leszy",
			Short: "SLAV command line interface",
			Long:  "SLAV command line interface, used to communicate with Boruta and Weles.",
			Run: func(cmd *cobra.Command, args []string) {
				err := cmd.Usage()
				if err != nil {
					logger.WithError(err).Error("Failed to print leszy command usage.")
				}
			},
		},
	},
	}
}

func (c *rootCmd) Cmd() *cobra.Command {
	return c.Command
}
