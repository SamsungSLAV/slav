package reqs

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/SamsungSLAV/slav/leszy"
	"github.com/SamsungSLAV/slav/logger"
)

//TODO: adjust marshalling of time.Time to more human readable form?

type ReqsListCmd struct {
	leszy.BaseCmd
	pretty bool
}

func NewReqsListCmd(c *leszy.Clients) *ReqsListCmd {
	rc := &ReqsListCmd{}
	return rc.New(c)
}

func (rc *ReqsListCmd) New(c *leszy.Clients) *ReqsListCmd {
	rc.Clients = c
	rc.Command = &cobra.Command{
		Use:   "list",
		Short: "List Boruta requests.",
		Long:  "", //TODO
		Run:   rc.ListReqs,
	}
	rc.addFlags()
	return rc
}

func (rc *ReqsListCmd) addFlags() {
	flagset := rc.Command.Flags()
	flagset.BoolVar(&rc.pretty, "pretty", true, "desc") //TODO: pretty should be global flag
}

func (rc *ReqsListCmd) Cmd() *cobra.Command {
	return rc.Command
}

func (rc *ReqsListCmd) ListReqs(cmd *cobra.Command, args []string) {
	rinfo, err := rc.Clients.Boruta.ListRequests(nil)
	if err != nil {
		logger.WithError(err).Error("Failed to list requests.")
	}
	var data []byte
	if data, err = json.Marshal(rinfo); err != nil {
		logger.WithError(err).Error("Failed to marshal Boruta's response.")
	}
	if rc.pretty {
		if data, err = json.MarshalIndent(rinfo, "", "    "); err != nil {
			logger.WithError(err).Error("Failed to marshal Boruta's reponse")
		}
	}
	fmt.Println(string(data))
}
