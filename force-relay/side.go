package main

import (
	"net"

	"github.com/eosforce/bus-service/force-relay/chainhandler"

	"github.com/cihub/seelog"
	commit "github.com/eosforce/bus-service/force-relay/pbs/relay"
	"github.com/eosforce/bus-service/force-relay/side"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startSideService() {
	lis, err := net.Listen("tcp", *transferURL)
	if err != nil {
		seelog.Errorf("failed to listen: %v", err)
		return
	}

	side.CreateClient(*configPath)
	relayCfg := side.NewCfg(*chain, *transfer)
	relayCfg.AppendActionInfo("force.token", "transfer")
	relayCfg.AppendActionInfo("force", "newaccount")
	side.SetCfg(relayCfg)

	service := grpc.NewServer()
	commit.RegisterRelayCommitServer(service,
		chainhandler.NewChainHandler(
			func(block *chainhandler.Block, actions []chainhandler.Action) {
				side.HandSideBlock(block, actions)
			}))
	reflection.Register(service)
	if err := service.Serve(lis); err != nil {
		seelog.Errorf("failed to serve: %v", err.Error())
	}
}
