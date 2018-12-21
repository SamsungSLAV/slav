package boruta

import (
	"github.com/spf13/cobra"

	"github.com/SamsungSLAV/slav/leszy"
	"github.com/SamsungSLAV/slav/logger"
)

type ReqsCmd struct {
	leszy.BaseCmd
}

func NewReqsCmd(c *leszy.Clients) *ReqsCmd {
	return &ReqsCmd{leszy.BaseCmd{
		Clients: c,
		Command: &cobra.Command{
			Use:   "reqs",
			Short: "Boruta requests management.",
			Long:  "", //TODO
			Run: func(cmd *cobra.Command, args []string) {
				err := cmd.Usage()
				if err != nil {
					logger.WithError(err).Error("Failed to print reqs command usage.")
				}
			},
		},
	}}

}

func (rc *ReqsCmd) Cmd() *cobra.Command {
	return rc.Command
}
