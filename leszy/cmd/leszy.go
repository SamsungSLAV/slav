package main

import (
	b "github.com/SamsungSLAV/boruta/http/client"
	"github.com/SamsungSLAV/slav/leszy"
	"github.com/SamsungSLAV/slav/leszy/boruta"
	"github.com/SamsungSLAV/slav/leszy/boruta/reqs"
	"github.com/SamsungSLAV/slav/leszy/boruta/workers"
)

func main() {
	// check leszy config
	// if empty- generate it and ask interactively for boruta and weles addr
	// integrate viper and bind viper output to cobra flags propagated through whole app

	clients := &leszy.Clients{
		Boruta: b.NewBorutaClient("http://106.120.47.240:8487"),
	}

	// define root path
	leszyCmd := leszy.NewRootCmd(clients)
	// create paths from root
	reqsCmd := boruta.NewReqsCmd(clients)
	workersCmd := boruta.NewWorkersCmd(clients)
	// register first level
	leszyCmd.Command.AddCommand(
		reqsCmd.Cmd(),
		workersCmd.Cmd(),
	)

	reqsListCmd := &reqs.ReqsListCmd{}
	reqsListCmd = reqsListCmd.New(clients)

	reqsNewCmd := &reqs.ReqsNewCmd{}
	reqsNewCmd = reqsNewCmd.New(clients)

	reqsCloseCmd := &reqs.ReqsCloseCmd{}
	reqsCloseCmd = reqsCloseCmd.New(clients)

	reqsAcquireCmd := &reqs.ReqsAcquireCmd{}
	reqsAcquireCmd = reqsAcquireCmd.New(clients)

	reqsCmd.Command.AddCommand(
		reqsListCmd.Command,
		reqsNewCmd.Command,
		reqsCloseCmd.Command,
		reqsAcquireCmd.Command)

	workersListCmd := &workers.WorkersListCmd{}
	workersListCmd = workersListCmd.New(clients)

	workersCmd.Command.AddCommand(workersListCmd.Command)

	leszyCmd.Command.Execute()
}
