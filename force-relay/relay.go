package main

import (
	"net"

	"github.com/eosforce/bus-service/force-relay/cfg"

	"github.com/eosforce/bus-service/force-relay/relay"

	"github.com/cihub/seelog"
	"github.com/eosforce/bus-service/force-relay/chainhandler"
	commit "github.com/eosforce/bus-service/force-relay/pbs/relay"
	"github.com/eosforce/bus-service/force-relay/side"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startRelayService() {
	lis, err := net.Listen("tcp", *relayURL)
	if err != nil {
		seelog.Errorf("failed to listen relay: %v", err)
		return
	}

	side.CreateClient(cfg.GetChainCfg("side"))
	relayCfg := side.NewCfg(*chain, *transfer)
	relayCfg.AppendActionInfo("relay.token", "destroy")
	side.SetCfg(relayCfg)

	service := grpc.NewServer()
	commit.RegisterRelayCommitServer(service,
		chainhandler.NewChainHandler(
			func(block *chainhandler.Block, actions []chainhandler.Action) {
				relay.HandRelayBlock(block, actions)
			}))
	reflection.Register(service)
	if err := service.Serve(lis); err != nil {
		seelog.Errorf("failed to serve: %v", err.Error())
	}
}
