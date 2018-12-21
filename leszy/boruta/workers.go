package boruta

import (
	"github.com/spf13/cobra"

	"github.com/SamsungSLAV/slav/leszy"
	"github.com/SamsungSLAV/slav/logger"
)

type WorkersCmd struct {
	leszy.BaseCmd
}

func NewWorkersCmd(c *leszy.Clients) *WorkersCmd {
	return &WorkersCmd{leszy.BaseCmd{
		Clients: c,
		Command: &cobra.Command{
			Use:   "workers",
			Short: "Boruta workers management.",
			Long:  "", //TODO
			Run: func(cmd *cobra.Command, args []string) {
				err := cmd.Usage()
				if err != nil {
					logger.WithError(err).Error("Failed to print workers command usage.")
				}
			},
		},
	}}

}

func (rc *WorkersCmd) Cmd() *cobra.Command {
	return rc.Command
}
