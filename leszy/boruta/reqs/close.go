package reqs

import (
	"strconv"

	"github.com/spf13/cobra"

	b "github.com/SamsungSLAV/boruta"
	"github.com/SamsungSLAV/slav/leszy"
	"github.com/SamsungSLAV/slav/logger"
)

type ReqsCloseCmd struct {
	leszy.BaseCmd
}

func NewReqsCloseCmd(c *leszy.Clients) *ReqsCloseCmd {
	rc := &ReqsCloseCmd{}
	return rc.New(c)
}

func (rc *ReqsCloseCmd) New(c *leszy.Clients) *ReqsCloseCmd {
	rc.Clients = c
	rc.Command = &cobra.Command{
		Use:   "close",
		Short: "Close Boruta request.",
		Long:  "", //TODO
		Args:  cobra.ExactArgs(1),
		Run:   rc.CloseReqs,
	}
	rc.addFlags()
	return rc
}

func (rc *ReqsCloseCmd) addFlags() {
	//flagset := rc.Command.Flags()
}

func (rc *ReqsCloseCmd) Cmd() *cobra.Command {
	return rc.Command
}

func (rc *ReqsCloseCmd) CloseReqs(cmd *cobra.Command, args []string) {
	rid, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		logger.WithError(err).Errorf("Failed to parse %s to uint64", args[0])
	}
	err = rc.Clients.Boruta.CloseRequest(b.ReqID(rid))
	if err != nil {
		logger.WithError(err).Error("Failed to create request.")
	}
}
