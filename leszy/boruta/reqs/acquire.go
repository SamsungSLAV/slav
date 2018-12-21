package reqs

import (
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"strconv"

	"golang.org/x/crypto/ssh"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	b "github.com/SamsungSLAV/boruta"
	"github.com/SamsungSLAV/slav/leszy"
	"github.com/SamsungSLAV/slav/logger"
)

type ReqsAcquireCmd struct {
	leszy.BaseCmd
}

func NewReqsAcquireCmd(c *leszy.Clients) *ReqsAcquireCmd {
	rc := &ReqsAcquireCmd{}
	return rc.New(c)
}

func (rc *ReqsAcquireCmd) New(c *leszy.Clients) *ReqsAcquireCmd {
	rc.Clients = c
	rc.Command = &cobra.Command{
		Use:   "acquire",
		Short: "When request is IN PROGRESS...", //TODO
		Long:  "",                               //TODO
		Args:  cobra.ExactArgs(1),
		Run:   rc.AcquireReqs,
	}
	rc.addFlags()
	return rc
}

func (rc *ReqsAcquireCmd) addFlags() {
	//flagset := rc.Command.Flags()
}

func (rc *ReqsAcquireCmd) Cmd() *cobra.Command {
	return rc.Command
}

func (rc *ReqsAcquireCmd) AcquireReqs(cmd *cobra.Command, args []string) {
	home, err := homedir.Dir()
	identityFilePath := home + "/.dryad-identity"

	rid, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		logger.WithError(err).Errorf("Failed to parse %s to uint64", args[0])
	}
	accessInfo, err := rc.Clients.Boruta.AcquireWorker(b.ReqID(rid))
	if err != nil {
		logger.WithError(err).Error("Failed to create request.")
	}
	//fmt.Println(accessInfo)
	//TODO: handle identity file in some more sensible way
	pub, err := ssh.NewPublicKey(&accessInfo.Key.PublicKey)
	key := ssh.MarshalAuthorizedKey(pub)

	err = ioutil.WriteFile(identityFilePath, key, 0777)
	if err != nil {
		logger.WithError(err).Error("Failed to create identity file.")
	}
	_, port, err := net.SplitHostPort(accessInfo.Addr.String())

	ccmd := exec.Command("/bin/bash", "-c", "ssh", "-G ",
		"-i "+identityFilePath, "-p "+port,
		accessInfo.Username+"@106.120.47.240")
	fmt.Println(ccmd)
	err = ccmd.Run()
	fmt.Println(err)
}
