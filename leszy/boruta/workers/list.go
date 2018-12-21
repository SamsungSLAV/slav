package workers

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/SamsungSLAV/slav/leszy"
	"github.com/SamsungSLAV/slav/logger"
)

type WorkersListCmd struct {
	leszy.BaseCmd
}

func NewWorkersListCmd(c *leszy.Clients) *WorkersListCmd {
	rc := &WorkersListCmd{}
	return rc.New(c)
}

func (rc *WorkersListCmd) New(c *leszy.Clients) *WorkersListCmd {
	rc.Clients = c
	rc.Command = &cobra.Command{
		Use:   "list",
		Short: "List Boruta workers.",
		Long:  "", //TODO
		Run:   rc.ListWorkers,
	}
	return rc
}

func (rc *WorkersListCmd) Cmd() *cobra.Command {
	return rc.Command
}

func (rc *WorkersListCmd) ListWorkers(cmd *cobra.Command, args []string) {
	rinfo, err := rc.Clients.Boruta.ListWorkers(nil, nil)
	if err != nil {
		logger.WithError(err).Error("Failed to list workers.")
	}
	fmt.Println(rinfo)
}
