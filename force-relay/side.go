package main

import (
	"net"

	"github.com/eosforce/bus-service/force-relay/cfg"
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

	side.InitCommitter()

	// grpc handler go
	go func() {
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
	}()

	// frome side need to commit block to relay
	side.CreateClient(cfg.GetChainCfg("relay"))
	side.StartCommitter()
}
